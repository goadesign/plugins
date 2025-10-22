package tools

import (
	"path/filepath"
	"sort"

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
	if len(g.toolsetOrder) == 0 {
		return nil
	}
	var files []*codegen.File
	for _, ts := range g.toolsetOrder {
		files = append(files, g.renderToolSet(ts)...)
	}
	return files
}

func (g *generator) renderToolSet(ts *toolSetData) []*codegen.File {
	var files []*codegen.File
	baseDir := toolSetOutputDir(ts)

	if pure := ts.pureTypes(); len(pure) > 0 {
		header := codegen.Header("Tool Types", ts.PackageName, ts.typeImports())
		files = append(files, &codegen.File{
			Path: filepath.Join(baseDir, "types.go"),
			SectionTemplates: []*codegen.SectionTemplate{
				header,
				{
					Name:   "tools-types",
					Source: toolsTemplates.Read(typesT),
					Data:   &typesTemplateData{Types: pure},
				},
			},
		})
	}

	if schemas := ts.schemaTypes(); len(schemas) > 0 {
		header := codegen.Header("Tool Schemas", ts.PackageName, nil)
		files = append(files, &codegen.File{
			Path: filepath.Join(baseDir, "schemas.go"),
			SectionTemplates: []*codegen.SectionTemplate{
				header,
				{
					Name:   "tools-schemas",
					Source: toolsTemplates.Read(schemasT),
					Data:   &schemasTemplateData{Types: schemas},
				},
			},
		})
	}

	if codecsFile := g.renderToolSetCodecs(ts, baseDir); codecsFile != nil {
		files = append(files, codecsFile)
	}

	if registryFile := g.renderToolSetRegistry(ts, baseDir); registryFile != nil {
		files = append(files, registryFile)
	}

	return files
}

func (g *generator) renderToolSetCodecs(ts *toolSetData, baseDir string) *codegen.File {
	if len(ts.Tools) == 0 {
		return nil
	}
	imports := []*codegen.ImportSpec{
		codegen.SimpleImport("encoding/json"),
		codegen.SimpleImport("fmt"),
		codegen.SimpleImport("goa.design/plugins/v3/tools"),
	}
	if ts.needsUnicodeImport() {
		imports = append(imports, codegen.SimpleImport("unicode/utf8"))
	}
	if ts.needsGoa() {
		imports = append(imports, codegen.GoaImport(""))
	}

	extra := make(map[string]*codegen.ImportSpec)
	for _, info := range ts.types() {
		if info.Import != nil && info.Import.Path != "" {
			extra[info.Import.Path] = info.Import
		}
		for _, im := range info.TypeImports {
			if im != nil && im.Path != "" {
				extra[im.Path] = im
			}
		}
	}
	if len(extra) > 0 {
		paths := make([]string, 0, len(extra))
		for p := range extra {
			paths = append(paths, p)
		}
		sort.Strings(paths)
		for _, p := range paths {
			imports = append(imports, extra[p])
		}
	}

	header := codegen.Header("Tool Codecs", ts.PackageName, imports)
	data := &codecsTemplateData{
		Types: ts.types(),
		Tools: ts.Tools,
	}
	return &codegen.File{
		Path: filepath.Join(baseDir, "codecs.go"),
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

func (g *generator) renderToolSetRegistry(ts *toolSetData, baseDir string) *codegen.File {
	if len(ts.Tools) == 0 {
		return nil
	}
	header := codegen.Header("Tool Registry", ts.PackageName, []*codegen.ImportSpec{
		codegen.SimpleImport("goa.design/plugins/v3/tools"),
	})
	return &codegen.File{
		Path: filepath.Join(baseDir, "registry.go"),
		SectionTemplates: []*codegen.SectionTemplate{
			header,
			{
				Name:   "tools-registry",
				Source: toolsTemplates.Read(registryT),
				Data:   &registryTemplateData{Tools: ts.Tools},
			},
		},
	}
}

func toolSetOutputDir(ts *toolSetData) string {
	if ts.ServicePath != "" {
		return filepath.Join(codegen.Gendir, ts.ServicePath, "tools", ts.DirName)
	}
	return filepath.Join(codegen.Gendir, "tools", ts.DirName)
}
