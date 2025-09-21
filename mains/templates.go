package mains

import (
    "embed"

    "goa.design/goa/v3/codegen/template"
)

//go:embed templates/*
var templatesFS embed.FS

var tmpl = &template.TemplateReader{FS: templatesFS}

