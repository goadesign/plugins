# Docs Plugin

The `docs` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that generates documentation from Goa designs. The plugin generates transport
agnostic documentation in the form of a JSON document structured as follows:

```
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
      "url": "https://goa.design/goa/v3"
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

To enable the plugin, import the docs DSL package in your design. Importing the DSL also registers the plugin automatically:

```go
import (
  . "goa.design/plugins/v3/docs/dsl"
  . "goa.design/goa/v3/dsl"
)
```

## Effects on Code Generation

Enabling the plugin changes the behavior of the `gen` command of the `goa` tool.
The command generates an additional `docs.json` at the top level containing the
documentation.

### Respecting openapi:generate metadata

The docs plugin respects Goa's `Meta("openapi:generate", "false")` (and legacy `"swagger:generate"`) to control documentation generation:

- Service-level: adding `Meta("openapi:generate", "false")` to a `Service` excludes that service from `docs.json`.
- Method-level: adding `Meta("openapi:generate", "false")` to a `Method` excludes that method from `docs.json`.
- HTTP-level: the same meta can be set inside `HTTP(...)` blocks for services or methods and will be honored.

This mirrors Goa's OpenAPI generation semantics.

### Disabling docs.json per service or method

To disable inclusion in `docs.json` without affecting OpenAPI generation, use the docs DSL `DisableDocs()` inside a service or method DSL:

```go
import (
  . "goa.design/plugins/v3/docs/dsl"
  . "goa.design/goa/v3/dsl"
)

var _ = Service("front", func() {
  DisableDocs() // service omitted from docs.json, OpenAPI still generated
  Method("about", func() { /* ... */ })
})

var _ = Service("S", func() {
  Method("M1", func() {
    DisableDocs() // method omitted from docs.json, OpenAPI still generated
    /* ... */
  })
})
```

Under the hood, `DisableDocs()` sets `Meta("docs:generate", "false")`; the plugin filters those out when building docs.json.

## Known Limitations

If `goa gen` is invoked with a custom output path (i.e. with the `-o` argument)
then the plugin appends to any pre-existing `doc.json` file instead of
overwriting. Make sure to delete the file prior to running `goa gen` when using
the `-o` option.

### Using JSON tags as field names

By default, the generated documentation uses attribute names as field names in definitions and examples. To instead use JSON struct tags declared via the `Meta("struct:tag:json", ...)` DSL, enable it in your design using the docs DSL (which also registers the plugin):

```go
import (
  . "goa.design/plugins/v3/docs/dsl"
  . "goa.design/goa/v3/dsl"
)

func init() {
  UseJSONTags()
}
```

When enabled, object property names in `definitions`, payloads, results, and error schemas/examples are renamed to match their JSON tag (the part before the first comma, e.g. `name,omitempty` becomes `name`). Fields tagged with `"-"` are ignored for renaming and keep their original names. The transformation does not mutate the Goa OpenAPI global definitions.

### Inlining $refs in JSON Schema

To inline `$ref` schemas where possible (while preserving cycles), enable it via the docs DSL:

```go
import (
  . "goa.design/plugins/v3/docs/dsl"
)

func init() {
  InlineRefs()
}
```
