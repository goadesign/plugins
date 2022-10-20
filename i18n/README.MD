# I18n Plugin

The `i18n` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that makes it possible to generate multiple internationalized openapi specs.
This plugin is decoupled from an i18n library, this is up to the user to decide. 

## Enabling the Plugin

To enable the plugin and make use of the I18n DSL simply import both the `i18n` and
the `dsl` packages as follows:

```go
import (
  i18n "goa.design/plugins/v3/i18n/dsl"
  . "goa.design/goa/v3/dsl"
)
```

## Setting up the environment

An environment variable `GOA_I18N` needs to be set up before running goa. This environment variable contains 
a comma separated list of all locales you wish to generate. For example: 

```
GOA_I18N=en goa gen calc/design          // Handles `en` locale only
GOA_I18N=en,nl goa gen calc/design       // Handles `en` and `nl` locales
GOA_I18N=en,nl,de_DE goa gen calc/design // Handles `en`, `nl` and `de_DE` locales
goa gen calc/design                      // Error, GOA_I18N not defined
```

The first locale in the list is the `default` locale. 

## Effects on Code Generation

Enabling the plugin changes the behavior of the `gen` command of the `goa` tool.
For each configured locale in `GOA_I18N` an additional `openapi_{locale}.(yaml|json)` is generated using the supplied locale.
The default locale (first in list `GOA_I18N`) will be used to generate `openapi.(yaml|json)` and will **not** generate an additional
locale specific spec. 


## Usage

The Goa `i18n` plugin is completely decoupled from the actual i18n implementation. This means you have to supply your own choice. 
These use case gives two examples how to do this. This first (Basic) showcases a simple key value store without any external dependencies. 
The second showcases an example using `go-18n` library (https://github.com/nicksnyder/go-i18n). 

### Basic 

```go
package design
// i18n.go
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

// design.go
package design 
import (
	. "goa.design/goa/v3/dsl"
	i18n "goa.design/plugins/v3/i18n/dsl"
)
var _ = API("calc",func(){
  i18n.Title(M("api.title"))
})
var _ = Service("calc", func() {  
  i18n.Description(M("calc.description"))
})

```

make sure to set the `GOA_I18N` environment variable and execute `Goa`: 

```
GOA_I18N=en,nl goa gen calc/design
```

### Advanced 

This advanced example uses `go-i18n` as implementation. 

```go
// i18n.go
import (
	"github.com/BurntSushi/toml"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func GetBundle() *i18n.Bundle {
	if bundle == nil {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		bundle.MustLoadMessageFile("i18n/messages.en.toml")
		bundle.MustLoadMessageFile("i18n/messages.nl.toml")
	}
	return bundle
}


func L(messageId string) func(locale string) string {
  return func(locale string) string {
		message, _ := i18n.NewLocalizer(GetBundle(), locale).Localize(messageConfig)
		return message
	}
}

// M is an alias for the Message struct
type M = i18n.Message
// LM is an alias for LocalizeMessage
func LM(m *M) func(locale string) string {
	return func(locale string) string {
		message, _ := i18n.NewLocalizer(GetBundle(), locale).LocalizeMessage(m)
		return message
	}
}

// design.go
import (
	. "goa.design/goa/v3/dsl"
	i18n "goa.design/plugins/v3/i18n/dsl"
)
var _ = API("calc",func(){
  i18n.Title(L("ApiCalcTitle"))
})
var _ = Service("calc", func() {  
  i18n.Description(LM(&M{
    ID: "ApiCalcServiceCalcDescription",
    Other: "Calculation service exposes mathematical operation endpoints"
  }))
})
```

make sure to set the `GOA_I18N` environment variable and execute `Goa`: 

```
GOA_I18N=en,nl goa gen calc/design
```

## Supported tags

Currently the following tags are supported: 
 * [Title](https://godoc.org/goa.design/goa/dsl#Title)
 * [Description](https://godoc.org/goa.design/goa/dsl#Description)
 * [Example](https://godoc.org/goa.design/goa/dsl#Example)
