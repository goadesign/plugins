package i18n

import (
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/i18n/expr"

	// Register code generators for the I18n plugin
	_ "goa.design/plugins/v3/i18n"
)

type Translateable = func(locale string) string

func Title(t Translateable) {
	i18n := &expr.I18nExpr{Trans: []Translateable{t}}

	current := eval.Current()
	i18n.Parent = current

	expr.Root.Title[current] = i18n
}

// Description adds a translatable description
func Description(t Translateable) {
	i18n := &expr.I18nExpr{Trans: []Translateable{t}}

	current := eval.Current()
	i18n.Parent = current

	expr.Root.Description[current] = i18n
}

// Example adds a translatable description
func Example(t ...Translateable) {
	i18n := &expr.I18nExpr{Trans: t}

	current := eval.Current()
	i18n.Parent = current

	expr.Root.Example[current] = i18n
}
