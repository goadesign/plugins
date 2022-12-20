package e2e

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

func init() {
	codegen.RegisterPlugin("e2e", "gen", nil, Generate)
}

func Example(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	return files, nil
}

func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	e2epath := filepath.Join(codegen.Gendir, "e2e")
	log.Printf("] %s [", e2epath)
	if _, err := os.Stat(e2epath); !os.IsNotExist(err) {
		// Goa does not delete files in the top-level gen folder.
		// https://github.com/goadesign/goa/pull/2194
		// The plugin must delete docs.json so that the generator does not append
		// to any existing docs.json.
		if err := os.Remove(e2epath); err != nil {
			panic(err)
		}
	}
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			files = append(files, e2eFiles(e2epath, r)...)
		}
	}
	return files, nil
}

func e2eFiles(e2epath string, r *expr.RootExpr) []*codegen.File {
	var retval []*codegen.File
	api := r.API
	for _, server := range api.Servers {
		for _, svc := range r.Services {
			if serviceIsOnServer(server.Services, svc.Name) {
				for _, method := range svc.Methods {
					f := &codegen.File{
						Path: fmt.Sprintf("%s/%s/%s/%s_test.go", e2epath, server.Name, svc.Name, method.Name),
					}
					title := fmt.Sprintf("%s/%s E2E Tests", svc.Name, method.Name)
					sections := []*codegen.SectionTemplate{
						codegen.Header(title, fmt.Sprintf("e2e_%s", server.Name), []*codegen.ImportSpec{
							{Path: "net/http"},
							{Path: "github.com/onsi/ginkgo/v2"},
							{Path: "github.com/onsi/gomega"},
						}),
					}
					sections = append(sections, &codegen.SectionTemplate{Name: "method-tests", Source: testMethodSection(), Data: method})
					f.SectionTemplates = sections
					retval = append(retval, f)
				}
			}
		}
		f := &codegen.File{
			Path: fmt.Sprintf("%s/%s/e2e_suite_test.go", e2epath, server.Name),
		}
		title := fmt.Sprintf("%s E2E Tests", api.Name)
		sections := []*codegen.SectionTemplate{
			codegen.Header(title, fmt.Sprintf("e2e_%s", server.Name), []*codegen.ImportSpec{
				{Path: "flag"},
				{Path: "github.com/onsi/ginkgo/v2"},
				{Path: "github.com/onsi/gomega"},
			}),
		}
		sections = append(sections, &codegen.SectionTemplate{Name: "flag-vars", Source: suiteParamVarsSection(), Data: nil})
		sections = append(sections, &codegen.SectionTemplate{Name: "flag-init", Source: suiteInitSection(), Data: nil})
		sections = append(sections, &codegen.SectionTemplate{Name: "suite-run", Source: suiteRunSection(), Data: nil})
		f.SectionTemplates = sections
		retval = append(retval, f)
	}
	for _, file := range retval {
		log.Printf("** %+v **", file.Path)
	}
	return retval
}

func serviceIsOnServer(services []string, service string) bool {
	for _, svc := range services {
		if svc == service {
			return true
		}
	}
	return false
}

func suiteParamVarsSection() string {
	return `{{ The vars to hold values from e2e command line params" | comment }}
	var baseUri`
}

func suiteInitSection() string {
	return `{{ Initialize flag params for command line e2e testing | comment }}
 func init() {
	flag.StringVar(&baseUri, "base-uri", "<<enter default host base uri here>>", "Base URI for tests to run against.")
 }`
}

func suiteRunSection() string {
	return `func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}`
}

func testMethodSection() string {
	return `var _ = Describe("{{ .Name }} ", Label("@{{ .Name }}"), func() {
	It("can test the method", func() {
		return true
	})
}`
}
