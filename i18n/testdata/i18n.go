package testdata

import (
	"encoding/json"
	"fmt"
)

// I18nLibrary is a mock library to showcase the functionality
type I18nLibrary struct {
	Messages map[string]map[string]string
}

// M returns a translated message for the specified key
func (i18nlib *I18nLibrary) M(label string) func(lang string) string {
	return func(lang string) string {
		messagesBundle, ok := i18nlib.Messages[lang]

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

var T = I18nLibrary{}
var _ = json.Unmarshal([]byte(`{
	"en": {
		"title": "Goa"
	}
}`), &T.Messages)
