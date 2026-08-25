package i18n_test

import (
	"reflect"
	"slices"
	"testing"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/i18n"
	"goa.design/plugins/v3/i18n/testdata"
)

func TestPrepare(t *testing.T) {
	cases := []struct {
		Name            string
		Locales         string
		DSL             func()
		ExpectedMessage string
	}{
		{"basic-usage", "en", testdata.SimpleI18nDSL, "Goa"},
		{"missing-locale", "nl", testdata.SimpleI18nDSL, "*title*"},
		{"default-locale", "en,nl", testdata.SimpleI18nDSL, "Goa"},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Setenv("GOA_I18N", c.Locales)
			root := codegen.RunDSL(t, c.DSL)

			roots := []eval.Root{root}

			// Expect Description to be empty
			checkExpr(roots, &expr.ServiceExpr{}, func(se interface{}) {
				if se.(*expr.ServiceExpr).Description != "" {
					t.Errorf("Description should be empty before prepare is run")
				}
			})
			i18n.Prepare("", roots)
			// Expect Description to be translated value
			checkExpr(roots, &expr.ServiceExpr{}, func(se interface{}) {
				d := se.(*expr.ServiceExpr).Description
				if d != c.ExpectedMessage {
					t.Errorf("Description %s does not match expected value %s", d, c.ExpectedMessage)
				}
			})
		})
	}
}
func checkExpr(roots []eval.Root, t interface{}, cb func(se interface{})) {
	for _, root := range roots {
		root.WalkSets(func(es eval.ExpressionSet) {
			for _, e := range es {
				et := reflect.TypeOf(e)
				if et == reflect.TypeOf(t) {
					cb(e)
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	t.Setenv("GOA_I18N", "en,nl")

	root := codegen.RunDSL(t, testdata.SimpleI18nDSL)
	roots := []eval.Root{root}
	i18n.Prepare("", roots)

	plan, err := httpcodegen.NewOpenAPIPlan(root, expr.NewExampleGenerator(root.API.RandomizerFactory))
	if err != nil {
		t.Fatal(err)
	}
	fs := plan.Files()
	gfs, _ := i18n.Generate("", roots, fs)

	if len(gfs) != 12 {
		t.Errorf("Expected to generate twelve files, received %d", len(gfs))
	}
	paths := make([]string, len(gfs))
	for i, f := range gfs {
		paths[i] = f.Path
	}
	for _, path := range []string{
		"gen/http/openapi3.2_nl.json",
		"gen/http/openapi3.2_nl.yaml",
	} {
		if !slices.Contains(paths, path) {
			t.Errorf("Expected generated files to contain %q, got %v", path, paths)
		}
	}
}

func TestGenerateWithOneLocalePreservesFiles(t *testing.T) {
	t.Setenv("GOA_I18N", "en")
	root := codegen.RunDSL(t, testdata.SimpleI18nDSL)
	files := []*codegen.File{{Path: "existing.go"}}

	generated, err := i18n.Generate("generated.local/gen", []eval.Root{root}, files)

	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(generated, files) {
		t.Fatalf("Generate returned %v, want the original files %v", generated, files)
	}
}
