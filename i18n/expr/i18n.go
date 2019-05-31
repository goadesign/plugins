package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
)

type Translateable = func(locale string) string
type (
	// I18nExpr describes a CORS policy.
	I18nExpr struct {
		// Origin is the origin string.
		Trans []Translateable
		// Parent expression, ServiceExpr or APIExpr.
		Parent eval.Expression
	}
)

func (i18n *I18nExpr) Messages(locale string) []string {
	messages := make([]string, len(i18n.Trans))
	for i, t := range i18n.Trans {
		messages[i] = t(locale)
	}
	return messages
}

// EvalName returns the generic expression name used in error messages.
func (i18n *I18nExpr) EvalName() string {
	var suffix string
	if i18n.Parent != nil {
		suffix = fmt.Sprintf(" of %s", i18n.Parent.EvalName())
	}
	return "I18N" + suffix
}

// Validate ensures the origin expression is valid.
func (i18n *I18nExpr) Validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	return verr
}
