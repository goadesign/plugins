package jaegertracer

import (
	"path/filepath"

	"goa.design/goa/codegen"
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
)

type fileToModify struct {
	file        *codegen.File
	path        string
	serviceName string
	isMain      bool
}

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginLast("jaegertracer-updater", "example", UpdateExample)
}

// UpdateExample modifies the example main file by adding
// the jaeger initialization section
func UpdateExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {

	filesToModify := []*fileToModify{}

	for _, root := range roots {
		if r, ok := root.(*httpdesign.RootExpr); ok {

			// Add the generated main files
			for _, svr := range r.Design.API.Servers {
				pkg := codegen.SnakeCase(codegen.Goify(svr.Name, true))
				mainPath := filepath.Join("cmd", pkg, "main.go")
				filesToModify = append(filesToModify, &fileToModify{path: mainPath, serviceName: svr.Name, isMain: true})
			}

			// Update the added files
			for _, fileToModify := range filesToModify {
				for _, file := range files {
					if file.Path == fileToModify.path {
						fileToModify.file = file
						updateExampleFile(genpkg, r, fileToModify)
						break
					}
				}
			}
		}
	}

	return files, nil
}

func updateExampleFile(genpkg string, root *httpdesign.RootExpr, f *fileToModify) {

	if f.isMain {

		sections := f.file.SectionTemplates
		header := sections[0]

		codegen.AddImport(header, &codegen.ImportSpec{
			Name: "jaeger",
			Path: "github.com/uber/jaeger-client-go",
		})

		codegen.AddImport(header, &codegen.ImportSpec{
			Name: "jaegercfg",
			Path: "github.com/uber/jaeger-client-go/config",
		})

		codegen.AddImport(header, &codegen.ImportSpec{
			Name: "jaegerlog",
			Path: "github.com/uber/jaeger-client-go/log",
		})

		index := -1
		var logSection *codegen.SectionTemplate

		for i, s := range sections {
			if s.Name == "service-main-logger" {
				index = i + 1
				logSection = s
				break
			}
		}

		if index > -1 {
			sections = append(sections, nil)
			copy(sections[index+1:], sections[index:])
			sections[index] = &codegen.SectionTemplate{
				Name:   "service-main-tracer",
				Source: tracerInitT,
				Data:   logSection.Data,
			}
			f.file.SectionTemplates = sections
		}
	}
}

const tracerInitT = `
	// Initialize a Jaeger tracer and assigns it as the global opentracing tracer
	tracerCfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: "{{ .TracingEndpoint }}",
		},
	}

	jLogger := jaegerlog.StdLogger

	closer, err := tracerCfg.InitGlobalTracer(
		"{{ .APIPkg }}",
		jaegercfg.Logger(jLogger),
	)

	if err != nil {
		log.Printf("could not initialize jaeger tracer: %s", err.Error())
	} else {
		defer closer.Close()
	}
`
