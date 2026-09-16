package i18n

import (
	"fmt"
	"os"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	goadsl "goa.design/goa/v3/dsl"
	"goa.design/goa/v3/eval"

	goaexpr "goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/goa/v3/http/codegen/openapi"
	"goa.design/plugins/v3/i18n/expr"
)

func init() {
	generator.RegisterPlugin("i18n", "gen", newPlugin)
}

// ENVKEY is the key used to lookup locales to use when producing translation openapi specs
const ENVKEY = "GOA_I18N"

func getLocales() ([]string, error) {
	locales := os.Getenv(ENVKEY)

	if locales == "" {
		return nil, fmt.Errorf("environment variable \"GOA_I18N\" not found, this is required to generate locale dependend output")
	}

	return strings.Split(locales, ","), nil
}

// Prepare executes all translations with the default language
func Prepare(_ string, roots []eval.Root) error {
	locales, error := getLocales()

	if error != nil {
		return error
	}

	defaultLocale := locales[0]
	walkTranslations(roots, defaultLocale)
	return nil
}

func walkTranslations(roots []eval.Root, locale string) {
	walker := func(s eval.ExpressionSet) {
		i18nRoot := expr.Root

		for _, e := range s {
			if i18nExpr, ok := i18nRoot.Description[e]; ok {
				handleDescriptionTranslation(e, i18nExpr.Messages(locale))
			}
			if i18nExpr, ok := i18nRoot.Example[e]; ok {
				handleExampleTranslation(e, i18nExpr.Messages(locale))
			}
			if i18nExpr, ok := i18nRoot.Title[e]; ok {
				handleTitleTranslation(e, i18nExpr.Messages(locale))
			}
		}
	}
	for _, root := range roots {
		root.WalkSets(walker)
	}
}

func handleDescriptionTranslation(p eval.Expression, d []string) {
	eval.Execute(func() {
		goadsl.Description(d[0])
	}, p)
}

func handleExampleTranslation(p eval.Expression, e []string) {
	eval.Execute(func() {
		goadsl.Example(e)
	}, p)
}
func handleTitleTranslation(p eval.Expression, e []string) {
	eval.Execute(func() {
		goadsl.Title(e[0])
	}, p)
}

// Generate produces additional OpenAPI files for locales configured through
// GOA_I18N. It remains available to callers that invoke the plugin directly.
func Generate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	plans, err := openAPIPlans(roots)
	if err != nil {
		return nil, err
	}
	if len(plans) <= 1 {
		return files, nil
	}
	for _, plan := range plans[1:] {
		files = append(files, plan.Files()...)
	}
	return files, nil
}

// newPlugin creates the state used by one Goa generation run.
func newPlugin() generator.Plugin {
	return generator.Plugin{
		Prepare: Prepare,
		Plan: func(plan *generator.Plan) error {
			plans, err := openAPIPlans(plan.Generation().Roots())
			if err != nil {
				return err
			}
			if len(plans) <= 1 {
				return nil
			}
			root := applicationRoot(plan.Generation().Roots())
			return plan.ReplaceOpenAPI(root, plans...)
		},
	}
}

// openAPIPlans builds the default OpenAPI documents followed by one copy for
// every additional locale. Each copy reads alternate values without changing
// the evaluated design used by other generators.
func openAPIPlans(roots []eval.Root) ([]*httpcodegen.OpenAPIPlan, error) {
	locales, err := getLocales()
	if err != nil {
		return nil, err
	}
	if len(locales) <= 1 {
		return nil, nil
	}
	root := applicationRoot(roots)
	examples := goaexpr.NewExampleGenerator(root.API.RandomizerFactory)
	base, err := httpcodegen.NewOpenAPIPlan(root, examples)
	if err != nil {
		return nil, err
	}
	plans := []*httpcodegen.OpenAPIPlan{base}
	for _, locale := range locales[1:] {
		specs, err := localizedSpecs(root, locale)
		if err != nil {
			return nil, err
		}
		localized, err := httpcodegen.NewOpenAPIPlanFromSpecs(
			root,
			goaexpr.NewExampleGenerator(root.API.RandomizerFactory),
			specs,
			localizedValues(roots, locale),
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, localized)
	}
	return plans, nil
}

// applicationRoot returns the Goa design root that owns the OpenAPI files.
func applicationRoot(roots []eval.Root) *goaexpr.RootExpr {
	for _, root := range roots {
		if result, ok := root.(*goaexpr.RootExpr); ok {
			return result
		}
	}
	panic("i18n: Goa design root is missing")
}

// localizedSpecs adds the locale before each OpenAPI filename extension.
func localizedSpecs(root *goaexpr.RootExpr, locale string) ([]openapi.Spec, error) {
	specs, err := openapi.Specs(root.API.Meta)
	if err != nil {
		return nil, err
	}
	for index := range specs {
		specs[index].Path += "_" + locale
	}
	return specs, nil
}

// localizedValues returns the translated text and examples for one OpenAPI
// document without changing the design read by other generators.
func localizedValues(roots []eval.Root, locale string) openapi.Values {
	values := openapi.Values{}
	for _, root := range roots {
		root.WalkSets(func(expressions eval.ExpressionSet) {
			for _, target := range expressions {
				if translation, ok := expr.Root.Description[target]; ok {
					values = values.WithDescription(target, translation.Messages(locale)[0])
				}
				if translation, ok := expr.Root.Title[target]; ok {
					values = values.WithTitle(target, translation.Messages(locale)[0])
				}
				if translation, ok := expr.Root.Example[target]; ok {
					if attribute, ok := target.(*goaexpr.AttributeExpr); ok {
						values = values.WithExamples(attribute, []*goaexpr.ExampleExpr{{
							Summary: "default",
							Value:   translation.Messages(locale),
						}})
					}
				}
			}
		})
	}
	return values
}
