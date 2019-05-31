package i18n_test

import (
	"os"
	"reflect"
	"testing"

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
			os.Setenv("GOA_I18N", c.Locales)
			httpcodegen.RunHTTPDSL(t, c.DSL)

			roots, _ := eval.Context.Roots()

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
		root.WalkSets(func(es eval.ExpressionSet) error {
			for _, e := range es {
				et := reflect.TypeOf(e)
				if et == reflect.TypeOf(t) {
					cb(e)
				}
			}
			return nil
		})
	}
}

func TestGenerate(t *testing.T) {
	os.Setenv("GOA_I18N", "en,nl")

	httpcodegen.RunHTTPDSL(t, testdata.SimpleI18nDSL)
	roots, _ := eval.Context.Roots()
	i18n.Prepare("", roots)

	fs, _ := httpcodegen.OpenAPIFiles(expr.Root)
	gfs, _ := i18n.Generate("", roots, fs)

	if len(gfs) != 4 {
		t.Errorf("Expected to generate four files, received %d", len(gfs))
	}
}
