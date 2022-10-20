package i18n

import (
	"fmt"
	"os"
	"strings"

	"goa.design/goa/v3/codegen"
	goadsl "goa.design/goa/v3/dsl"
	"goa.design/goa/v3/eval"

	goaexpr "goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/i18n/expr"
)

func init() {
	codegen.RegisterPlugin("i18n", "gen", Prepare, Generate)
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
func Prepare(genpkg string, roots []eval.Root) error {
	locales, error := getLocales()

	if error != nil {
		return error
	}

	defaultLocale := locales[0]
	walkTranslations(roots, defaultLocale)
	return nil
}

func walkTranslations(roots []eval.Root, locale string) {
	walker := func(s eval.ExpressionSet) error {
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
		return nil
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

// Generate produces additional openapi files for locales configured via
// the system environment variable GOA_I18N
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	locales, _ := getLocales()

	if len(locales) <= 1 {
		// Nothing to generate, default already contains translations of default locale
		return files, nil
	}
	defaultLocale := locales[0]
	restLocales := locales[1:]

	for _, locale := range restLocales {
		walkTranslations(roots, locale)

		fs, _ := httpcodegen.OpenAPIFiles(goaexpr.Root)
		// Rename the files
		for _, file := range fs {
			sp := strings.Split(file.Path, ".")
			file.Path = fmt.Sprintf("%s_%s.%s", sp[0], locale, sp[1])
		}
		files = append(files, fs...)
	}
	// Not sure if needed, but reset messages to defaultLocale
	walkTranslations(roots, defaultLocale)
	return files, nil
}
