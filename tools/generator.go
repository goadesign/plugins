package tools

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
)

func init() {
	codegen.RegisterPlugin("tools", "gen", nil, Generate)
}

// Generate emits the tool metadata, schemas, and codec files under per-toolset
// packages scoped by service (gen/<service>/tools/<toolset>/).
func Generate(genpkg string, _ []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	g, err := newGenerator(genpkg)
	if err != nil {
		return nil, err
	}
	if err := g.collect(); err != nil {
		return nil, err
	}
	out := g.render()
	if len(out) == 0 {
		return files, nil
	}
	return append(files, out...), nil
}
