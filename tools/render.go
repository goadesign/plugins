package tools

import (
	"path/filepath"
	"sort"
	"strings"

	"goa.design/goa/v3/codegen"
)

type (
	typesTemplateData struct {
		Types []*typeInfo
	}
	schemasTemplateData struct {
		Types []*typeInfo
	}
	codecsTemplateData struct {
		Types []*typeInfo
		Tools []*toolEntry
	}
	registryTemplateData struct {
		Tools []*toolEntry
	}
)

func (g *generator) render() []*codegen.File {
	if len(g.tools) == 0 {
		return nil
	}
	var files []*codegen.File
	if f := g.renderTypes(); f != nil {
		files = append(files, f)
	}
	if f := g.renderSchemas(); f != nil {
		files = append(files, f)
	}
	if f := g.renderCodecs(); f != nil {
		files = append(files, f)
	}
	if f := g.renderRegistry(); f != nil {
		files = append(files, f)
	}
	return files
}

func (g *generator) renderTypes() *codegen.File {
	var pure []*typeInfo
	for _, info := range g.ordered {
		if info.NeedType {
			pure = append(pure, info)
		}
	}
	if len(pure) == 0 {
		return nil
	}
	header := codegen.Header("Tool Types", "tools", g.typeImportList())
	return &codegen.File{
		Path: filepath.Join(codegen.Gendir, "tools", "types.go"),
		SectionTemplates: []*codegen.SectionTemplate{
			header,
			{
				Name:   "tools-types",
				Source: toolsTemplates.Read(typesT),
				Data:   &typesTemplateData{Types: pure},
			},
		},
	}
}

func (g *generator) renderSchemas() *codegen.File {
	var schemas []*typeInfo
	for _, info := range g.ordered {
		if info.SchemaLiteral != "" {
			schemas = append(schemas, info)
		}
	}
	if len(schemas) == 0 {
		return nil
	}
	header := codegen.Header("Tool Schemas", "tools", nil)
	return &codegen.File{
		Path: filepath.Join(codegen.Gendir, "tools", "schemas.go"),
		SectionTemplates: []*codegen.SectionTemplate{
			header,
			{
				Name:   "tools-schemas",
				Source: toolsTemplates.Read(schemasT),
				Data:   &schemasTemplateData{Types: schemas},
			},
		},
	}
}

func (g *generator) renderCodecs() *codegen.File {
	imports := []*codegen.ImportSpec{
		codegen.SimpleImport("encoding/json"),
		codegen.SimpleImport("fmt"),
		codegen.SimpleImport("goa.design/plugins/v3/tools"),
	}
	if g.needsUnicodeImport() {
		imports = append(imports, codegen.SimpleImport("unicode/utf8"))
	}
	if g.needsGoa {
		imports = append(imports, codegen.GoaImport(""))
	}
	if len(g.imports) > 0 {
		paths := make([]string, 0, len(g.imports))
		for p := range g.imports {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			imports = append(imports, g.imports[p])
		}
	}
	header := codegen.Header("Tool Codecs", "tools", imports)
	data := &codecsTemplateData{
		Types: g.ordered,
		Tools: g.tools,
	}
	return &codegen.File{
		Path: filepath.Join(codegen.Gendir, "tools", "codecs.go"),
		SectionTemplates: []*codegen.SectionTemplate{
			header,
			{
				Name:   "tools-codecs",
				Source: toolsTemplates.Read(codecsT),
				Data:   data,
			},
		},
	}
}

func (g *generator) renderRegistry() *codegen.File {
	if len(g.tools) == 0 {
		return nil
	}
	header := codegen.Header("Tool Registry", "tools", []*codegen.ImportSpec{
		codegen.SimpleImport("goa.design/plugins/v3/tools"),
	})
	return &codegen.File{
		Path: filepath.Join(codegen.Gendir, "tools", "registry.go"),
		SectionTemplates: []*codegen.SectionTemplate{
			header,
			{
				Name:   "tools-registry",
				Source: toolsTemplates.Read(registryT),
				Data:   &registryTemplateData{Tools: g.tools},
			},
		},
	}
}

func (g *generator) typeImportList() []*codegen.ImportSpec {
	if len(g.typeImports) == 0 {
		return nil
	}
	paths := make([]string, 0, len(g.typeImports))
	for p := range g.typeImports {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	imports := make([]*codegen.ImportSpec, 0, len(paths))
	for _, p := range paths {
		imports = append(imports, g.typeImports[p])
	}
	return imports
}

func (g *generator) needsUnicodeImport() bool {
	for _, info := range g.ordered {
		if info.HasValidation && strings.Contains(info.Validation, "utf8.") {
			return true
		}
	}
	return false
}
