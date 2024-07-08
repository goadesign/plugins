package zerologger

import (
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPluginLast("zerologger", "example", nil, UpdateExample)
}

// UpdateExample modifies the example generated files by replacing
// the log import reference when needed
// It also modify the initially generated main and service files
func UpdateExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		updateExample(f)
	}
	return files, nil
}

func updateExample(file *codegen.File) {
	for _, section := range file.SectionTemplates {
		switch section.Name {
		case "server-main-services":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/rs/zerolog"})
			oldinit := "{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}()"
			section.Source = strings.Replace(section.Source, oldinit, initT, 1)
		case "basic-service-struct":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "github.com/rs/zerolog"})
			section.Source = basicServiceStructT
		case "basic-service-init":
			section.Source = basicServiceInitT
		case "basic-endpoint":
			section.Source = strings.Replace(
				section.Source,
				`log.Printf(ctx, "{{ .ServiceVarName }}.{{ .Name }}")`,
				`s.logger.Info().Msgf("{{ .ServiceVarName}}.{{ .Name }}")`,
				1,
			)
		}
	}
}

const (
	initT = `logLevel := zerolog.InfoLevel
	if *dbgF {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Str("service", serviceName).Logger()
	{{ .VarName }}Svc = {{ $.APIPkg }}.New{{ .StructName }}(logger)`

	basicServiceInitT = `
{{ printf "New%s returns the %s service implementation." .StructName .Name | comment }}
func New{{ .StructName }}(logger *zerolog.Logger) {{ .PkgName }}.Service {
	return &{{ .VarName }}srvc{
		logger: logger,
	}
}
`

	basicServiceStructT = `
	{{ printf "%s service example implementation.\nThe example methods log the requests and return zero values." .Name | comment }}
	type {{ .VarName }}srvc struct {
		logger *zerolog.Logger
	}
`
)
