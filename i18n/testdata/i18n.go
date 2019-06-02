package testdata

import (
	"encoding/json"
	"fmt"
)

// Messages is a nested key value store where
// localized messages are store in format <locale>:<label>=<translation>
var messages map[string]map[string]string
var _ = json.Unmarshal([]byte(`{
	"en": {
		"title": "Goa"
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
