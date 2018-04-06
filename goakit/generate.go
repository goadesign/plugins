package goakit

import (
	"fmt"
	"regexp"

	"goa.design/goa/codegen"
	"goa.design/goa/eval"
	httpdesign "goa.design/goa/http/design"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginFirst("goakit", "gen", Generate)
	codegen.RegisterPluginFirst("goakit", "example", Example)
	codegen.RegisterPluginLast("goakit-goakitify", "gen", Goakitify)
}

// Generate generates go-kit specific decoders and encoders.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*httpdesign.RootExpr); ok {
			files = append(files, EncodeDecodeFiles(genpkg, r)...)
			files = append(files, MountFiles(r)...)
		}
	}
	return files, nil
}

// Goakitify modifies all the previously generated files by replacing all
// instances of "goa.Endpoint" with "github.com/go-kit/kit/endpoint".Endpoint
// and adding the corresponding import.
func Goakitify(enpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		goakitify(f)
	}
	return files, nil
}

// Example iterates through the roots and returns files that implement an
// example service and client.
func Example(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	var examples []*codegen.File
	for _, root := range roots {
		if r, ok := root.(*httpdesign.RootExpr); ok {
			examples = ExampleServerFiles(genpkg, r)
			break
		}
	}
	if len(examples) == 0 {
		return nil, fmt.Errorf("example: no HTTP design found")
	}
	// Remove previously generated example files.
	var output []*codegen.File
	for _, f := range files {
		found := false
		for _, ex := range examples {
			if f.Path == ex.Path {
				found = true
				break
			}
		}
		if !found {
			output = append(output, f)
		}
	}
	output = append(output, examples...)
	return output, nil
}

// goaEndpointRegexp matches occurrences of the "goa.Endpoint" type in Go code.
var goaEndpointRegexp = regexp.MustCompile(`([^\p{L}_])goa\.Endpoint([^\p{L}_])`)

// goakitify replaces all occurrences of goa.Endpoint with endpoint.Endpoint in
// the file section template sources. It also adds
// "github.com/go-kit/kit/endpoint" to the list of imported packages if
// occurrences were replaces.
func goakitify(f *codegen.File) {
	var hasEndpoint bool
	for _, s := range f.SectionTemplates {
		if !hasEndpoint {
			hasEndpoint = goaEndpointRegexp.MatchString(s.Source)
		}
		s.Source = goaEndpointRegexp.ReplaceAllString(s.Source, "${1}endpoint.Endpoint${2}")
	}
	if hasEndpoint {
		codegen.AddImport(
			f.SectionTemplates[0],
			&codegen.ImportSpec{Path: "github.com/go-kit/kit/endpoint"},
		)
	}
}
