package tools

import (
	"embed"

	"goa.design/goa/v3/codegen/template"
)

const (
	typesT    = "types"
	schemasT  = "schemas"
	codecsT   = "codecs"
	registryT = "registry"
)

//go:embed templates/*
var templateFS embed.FS

var toolsTemplates = &template.TemplateReader{FS: templateFS}
