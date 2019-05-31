# I18n Plugin

The `i18n` plugin is a [goa v3](https://github.com/goadesign/goa/tree/v3) plugin
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

## Effects on Code Generation

Enabling the plugin changes the behavior of both the `gen` command of the `goa` tool.


## Usage

### Basic 



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

// MC is an alias for the LocalizeConfig
type MC = i18n.LocalizeConfig

// M returns the translated message according to the message config supplied
func M(messageConfig *MC) func(locale string) string {
	return func(locale string) string {
		message, _ := i18n.NewLocalizer(GetBundle(), locale, "").Localize(messageConfig)
		return message
	}
}

// design.go
import (
	. "goa.design/goa/v3/dsl"
	i18n "goa.design/plugins/v3/i18n/dsl"
)
var _ = API("calc",func(){
  i18n.Title(M(&MC{
    MessageID: "api.title" 
  }))
})
var _ = Service("calc", func() {  
  i18n.Description(M(&MC{
    MessageID: "calc.description" 
  }))
})
```

make sure to set the `GOA_I18N` environment variable and execute `Goa`: 

```
GOA_I18N=en,nl goa gen calc/design
```