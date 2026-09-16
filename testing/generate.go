// This file registers the testing plugin and supports both current per-run Goa
// plans and the released callback API used by existing applications.
package testing

import (
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	testcodegen "goa.design/plugins/v3/testing/codegen"
)

// Generate produces test harness files for each service in the design (gen phase).
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	services, err := plannedServices(genpkg, roots)
	if err != nil {
		return nil, err
	}

	for _, root := range roots {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		for _, svc := range r.Services {
			if sd := services[r].Get(svc.Name); sd != nil {
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
	services, err := plannedServices(genpkg, roots)
	if err != nil {
		return nil, err
	}
	for _, root := range roots {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		scope := codegen.NewNameScope()
		for _, svc := range r.Services {
			s := services[r].Get(svc.Name)
			if s == nil {
				continue
			}
			scope.Unique(s.PkgName)
		}
		apipkg := scope.Unique(strings.ToLower(codegen.Goify(r.API.Name, false)), "api")
		for _, svc := range r.Services {
			svcData := services[r].Get(svc.Name)
			if svcData == nil {
				continue
			}
			if f := testcodegen.GenerateSuiteTopLevelFromData(genpkg, apipkg, svcData, r, svc); f != nil {
				files = append(files, f)
			}
			// Generate example scenarios.yaml (only for first service to avoid duplicates)
			if f := testcodegen.GenerateExampleScenariosFromData(genpkg, svcData, r, svc); f != nil {
				files = append(files, f)
				break // Only generate one scenarios.yaml for all services
			}
		}
	}
	return files, nil
}

// init registers fresh plugin callbacks for every Goa generation run.
func init() {
	generator.RegisterPlugin("testing", "gen", newGeneratePlugin)
	generator.RegisterPlugin("testing", "example", newExamplePlugin)
}

// newGeneratePlugin creates the callback that reads finalized Goa plans for
// one generation run and appends generated testing helpers.
func newGeneratePlugin() generator.Plugin {
	return generator.Plugin{
		Plan:     planGeneratedPackages,
		Generate: generatePlanned,
	}
}

// newExamplePlugin creates the callback that reads finalized Goa plans for
// one example run and appends user-owned testing files when they do not exist.
func newExamplePlugin() generator.Plugin {
	return generator.Plugin{Generate: generatePlannedExample}
}

// planGeneratedPackages records the imports and fixed public names emitted by
// each generated service testing package.
func planGeneratedPackages(plan *generator.Plan) error {
	for _, candidate := range plan.Generation().Roots() {
		root, ok := candidate.(*expr.RootExpr)
		if !ok {
			continue
		}
		servicePlan := plan.Service(root)
		for _, svc := range root.Services {
			if err := testcodegen.PlanPackage(plan.Generation(), servicePlan, root, svc); err != nil {
				return err
			}
		}
	}
	return nil
}

// generatePlanned appends testing helpers using the service names and
// JSON-RPC endpoints finalized in the current Goa generation plan.
func generatePlanned(plan *generator.Plan, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range plan.Generation().Roots() {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		services := plan.Service(r).Services()
		jsonPlan, _ := plan.JSONRPC(r)
		for _, svc := range r.Services {
			if generated := testcodegen.GeneratePlanned(plan.Generation().GenPkg(), services.Get(svc.Name), r, svc, jsonPlan); len(generated) > 0 {
				files = append(files, generated...)
			}
		}
	}
	return files, nil
}

// generatePlannedExample appends the user-owned suite and scenario files using
// the names and JSON-RPC endpoints finalized in the current Goa plan.
func generatePlannedExample(plan *generator.Plan, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range plan.Generation().Roots() {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		services := plan.Service(r).Services()
		scope := codegen.NewNameScope()
		for _, svc := range r.Services {
			scope.Unique(services.Get(svc.Name).PkgName)
		}
		apiPackage := scope.Unique(strings.ToLower(codegen.Goify(r.API.Name, false)), "api")
		jsonPlan, _ := plan.JSONRPC(r)
		for _, svc := range r.Services {
			svcData := services.Get(svc.Name)
			if file := testcodegen.GeneratePlannedSuiteTopLevel(plan.Generation().GenPkg(), apiPackage, svcData, r, svc, jsonPlan); file != nil {
				files = append(files, file)
			}
			if file := testcodegen.GeneratePlannedExampleScenarios(plan.Generation().GenPkg(), svcData, r, svc, jsonPlan); file != nil {
				files = append(files, file)
				break
			}
		}
	}
	return files, nil
}

// plannedServices runs Goa's current service planner for released callbacks,
// which receive roots and files but not the generation plan used by Goa.
func plannedServices(genpkg string, roots []eval.Root) (map[*expr.RootExpr]*service.ServicesData, error) {
	generation, err := codegen.NewGeneration(genpkg, roots)
	if err != nil {
		return nil, err
	}
	inputs := make([]service.PlanInput, 0, len(roots))
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			inputs = append(inputs, service.PlanInput{
				Root:     r,
				Examples: expr.NewExampleGenerator(r.API.RandomizerFactory),
			})
		}
	}
	plans, err := service.NewPlans(generation, inputs...)
	if err != nil {
		return nil, err
	}
	if err := generation.Freeze(); err != nil {
		return nil, err
	}
	services := make(map[*expr.RootExpr]*service.ServicesData, len(plans))
	for _, plan := range plans {
		if err := plan.Link(); err != nil {
			return nil, err
		}
		services[plan.Root()] = plan.Services()
	}
	return services, nil
}
