// This file changes Goa's starter service and server examples to create and
// pass a Zap logger. The plugin updates both sides of the constructor call so
// every generated example remains ready to build.
package zaplogger

import (
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

// init registers the example update after Goa has built its starter files.
func init() {
	codegen.RegisterPluginLast("zaplogger", "example", nil, UpdateExample)
}

// UpdateExample replaces the starter logger in each generated service and
// passes the matching Zap logger from the server main.
func UpdateExample(_ string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, f := range files {
		updateExample(f)
	}
	return files, nil
}

// updateExample replaces the matching sections in one generated starter file.
func updateExample(file *codegen.File) {
	for _, section := range file.SectionTemplates {
		switch section.Name {
		case "server-main-services":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "go.uber.org/zap"})
			oldinit := "{{ .ServiceVar }} = {{ $.APIPkg }}.{{ .ExampleConstructorDeclaration.Name }}()"
			section.Source = strings.Replace(section.Source, oldinit, initT, 1)
		case "basic-service-struct":
			codegen.AddImport(file.SectionTemplates[0], &codegen.ImportSpec{Path: "go.uber.org/zap"})
			section.Source = basicServiceStructT
		case "basic-service-init":
			section.Source = basicServiceInitT
		case "basic-endpoint":
			section.Source = strings.Replace(
				section.Source,
				`log.Printf(ctx, "{{ .ServiceVarName }}.{{ .Name }}")`,
				`s.logger.Info("{{ .ServiceVarName}}.{{ .Name }}")`,
				1,
			)
		}
	}
}

const (
	initT = `
	var zlog *zap.SugaredLogger
	if *dbgF {
		l, _ := zap.NewDevelopment()
		zlog = l.Sugar().With(zap.String("service", {{ printf "%q" .Name }}))
	} else {
		l, _ := zap.NewProduction()
		zlog = l.Sugar().With(zap.String("service", {{ printf "%q" .Name }}))
	}
	{{ .ServiceVar }} = {{ $.APIPkg }}.{{ .ExampleConstructorDeclaration.Name }}(zlog)`

	basicServiceInitT = `
{{ printf "New%s returns the %s service implementation." .StructName .Name | comment }}
func New{{ .StructName }}(logger *zap.SugaredLogger) {{ .PkgName }}.Service {
	return &{{ .VarName }}srvc{
		logger: logger,
	}
}
`

	basicServiceStructT = `
	{{ printf "%s service example implementation.\nThe example methods log the requests and return zero values." .Name | comment }}
	type {{ .VarName }}srvc struct {
		logger *zap.SugaredLogger
	}
`
)
