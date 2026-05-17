package testing

import (
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	testcodegen "goa.design/plugins/v3/testing/codegen"
)

// Register registers the testing plugin with Goa.
func init() {
	codegen.RegisterPlugin("testing", "gen", nil, Generate)
	codegen.RegisterPlugin("testing", "example", nil, GenerateExample)
}

// Generate produces test harness files for each service in the design (gen phase).
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	// Collect service.Data from incoming files
	svcByName := map[string]*service.Data{}
	for _, f := range files {
		for _, s := range f.SectionTemplates {
			if sd, ok := s.Data.(*service.Data); ok && sd != nil {
				svcByName[sd.Name] = sd
			}
		}
	}

	for _, root := range roots {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		for _, svc := range r.Services {
			if sd := svcByName[svc.Name]; sd != nil {
				if fs := testcodegen.Generate(genpkg, sd, r, svc); len(fs) > 0 {
					files = append(files, fs...)
				}
			}
		}
	}
	return files, nil
}

// GenerateExample produces the top-level user-editable test suite files using only DSL data (example phase).
func GenerateExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		for _, svc := range r.Services {
			// Get example implementation package name from actual file header
			var apipkg string
			for _, f := range files {
				for _, section := range f.Section("source-header") {
					header, ok := section.Data.(map[string]any)
					if !ok {
						continue
					}
					p, ok := header["Pkg"].(string)
					if !ok {
						continue
					}
					apipkg = p
					break
				}
				if apipkg != "" {
					break
				}
			}
			if f := testcodegen.GenerateSuiteTopLevel(genpkg, apipkg, r, svc); f != nil {
				files = append(files, f)
			}
			// Generate example scenarios.yaml (only for first service to avoid duplicates)
			if f := testcodegen.GenerateExampleScenarios(genpkg, r, svc); f != nil {
				files = append(files, f)
				break // Only generate one scenarios.yaml for all services
			}
		}
	}
	return files, nil
}
