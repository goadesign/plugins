package design

import (
	"encoding/json"
	"fmt"
)

// Messages is a nested key value store where
// localized messages are store in format <locale>:<label>=<translation>
var messages map[string]map[string]string
var _ = json.Unmarshal([]byte(`{
	"en": {
		"ApiCalcTitle":"CORS Example Calc API",
		"ApiCalcDescription":"This API demonstrates the use of the goa I18n plugin",
		"ApiCalcServiceCalcDescription":"The calc service exposes public endpoints to do basic mathematical calculations.",
		"ApiCalcServiceCalcMethodAddDescription":"Add adds up the two integer parameters and returns the results.",
		"ApiCalcServiceCalcMethodAddPayloadAttributeADescription":"Left operand",
		"ApiCalcServiceCalcMethodAddPayloadAttributeBDescription":"Right operand",
		"ApiCalcServiceCalcMethodAddResultDescription":"Result of addition"
	},
	"nl": {
		"ApiCalcTitle":"CORS Voorbeeld Calc API",
		"ApiCalcDescription":"Dit is een demonstratie van de vertalings plugin (i18n) van Goa",
		"ApiCalcServiceCalcDescription":"De reken service stelt basis rekenmethodes publiekelijk beschikbaar",
		"ApiCalcServiceCalcMethodAddDescription":"Tel twee getallen bij elkaar op en retourneerd het resultaat.",
		"ApiCalcServiceCalcMethodAddPayloadAttributeADescription":"Linker operand",
		"ApiCalcServiceCalcMethodAddPayloadAttributeBDescription":"Rechter operand",
		"ApiCalcServiceCalcMethodAddResultDescription":"Resultaat van optellen"
	}
}`), &messages)

// M returns a translated message for the specified key
func M(label string) func(lang string) string {
	return func(lang string) string {
		messagesBundle, ok := messages[lang]

		if !ok {
			return fmt.Sprintf("*%s*", label)
		}
		message, ok := messagesBundle[label]
		if !ok {
			return fmt.Sprintf("*%s*", label)
		}
		return message
	}
}
