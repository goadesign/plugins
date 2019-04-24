# Docs Plugin

The `docs` plugin is a [Goa](https://github.com/goadesign/goa/tree/v2) plugin
that generates documentation from Goa designs. The plugin generates transport
agnostic documentation in the form of a JSON document structured as follows:

```json
{
  "api": {
    "name": "API A",
    "title": "an API",
    "description": "an API",
    "version": "v1",
    "servers": {
      "server A": {
        "name": "server A",
        "description": "a server",
        "services": ["service A"],
        "hosts": {
          "dev": {
            "name": "dev",
            "server": "iis",
            "uris": ["http://localhost:80", "https://localhost:80"]
          }
        }
      }
    },
    "terms": "the terms of the API",
    "contact": {
      "name": "support",
      "email": "support@goa.design"
    },
    "license": {
      "name": "MIT",
    },
    "docs": {
      "url": "https://goa.design/goa"
    },
    "requirements": [{
      "schemes": [{
        "name": "scheme A",
        "type": "jwt"
      }],
      "scopes": ["api:read"]
    }]
  },
  "services": {
    "service A": {
      "name": "service A",
      "description": "a service",
      "methods": {
        "method A": {
          "name": "method A",
          "description": "a method",
          "payload": {
            "type": "#/types/typeA",
            "streaming": false,
            "example": { /* a valid instance of typeA */ },
          },
          "result": {
            "type": "#/types/typeB",
            "streaming": false,
            "example": { /* a valid instance of typeB */ },
          },
          "errors": {
            "error A": {
              "name": "error A",
              "description": "an error",
              "type": "#/types/typeC",
              "example": { /* a valid instance of typeC */ },
            },
            "requirements": [{
              "schemes": [{
                "name": "scheme A",
                "type": "jwt"
              }],
              "scopes": ["api:read"]
            }],
          },
          "requirements": [{
            "schemes": [{
              "name": "scheme A",
              "type": "jwt"
            }],
            "scopes": ["api:read"]
          }],
        }
      }
    }
  },
  "types": {
    "typeA": { /* JSON schema describing type A */ },
    "typeB": { /* JSON schema describing type B */ },
    "typeC": { /* JSON schema describing type C */ }
  }
}
```

## Enabling the Plugin

To enable the plugin simply import both the `docs` package as follows:

```go
import (
  _ "goa.design/plugins/docs"
  . "goa.design/goa/dsl"
)
```
Note the use of blank identifier to import the `docs` package which is necessary
as the package is imported solely for its side-effects (initialization).

## Effects on Code Generation

Enabling the plugin changes the behavior of the `gen` command of the `goa` tool.
The command generates an additional `doc.json` at the top level containing the 
documentation.
