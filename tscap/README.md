# Tailscale App Capabilities Plugin

The `tscap` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that provides declarative authorization using Tailscale app capabilities.

Inspired by the implementation of [arnz](../arnz/README.md).

## Requirements

- Tailscale 1.92+ (app capabilities feature)
- Service must be served via `tailscale serve --accept-app-caps`

## Enabling the Plugin

To enable the plugin and make use of the tscap DSL simply import both the `tscap`
and the `dsl` packages as follows:

```go
import (
  . "goa.design/goa/v3/dsl"
  tscap "goa.design/plugins/v3/tscap/dsl"
)
```

### Tailscale Setup

```bash
tailscale serve --accept-app-caps example.com/cap/myapp https+insecure://localhost:8080
```

### ACL Grants

Configure grants in your tailnet policy:

```json
{
  "grants": [
    {
      "src": ["group:developers"],
      "dst": ["tag:myapp"],
      "app": {
        "example.com/cap/myapp": [{"action": ["*"], "resources": ["*"]}]
      }
    },
    {
      "src": ["group:finance"],
      "dst": ["tag:myapp"],
      "app": {
        "example.com/cap/myapp": [{"action": ["read"], "resources": ["items/*"]}]
      }
    }
  ]
}
```

## Effects on Code Generation

Enabling the plugin changes the behavior of the `gen` command of the `goa` tool.

The `gen` command output is modified as follows:

1. Generates middleware that extracts the `Tailscale-App-Capabilities` header
2. Parses the JSON capabilities from the header
3. Checks if the caller's grants satisfy the method's requirements
4. Returns 401 if header is missing, 403 if permissions are insufficient

## Design

This plugin adds the following functions to the Goa DSL:

* `Require` declares that the method requires a Tailscale app capability with the
  specified action and resource.
* `AllowAnonymous` marks the method as not requiring any capability check. Requests
  without the capabilities header will be allowed through.

The usage and effect of the DSL functions are described in the [Godocs](https://godoc.org/goa.design/plugins/v3/tscap/dsl)

Here is an example defining capability requirements at a method level.

```go
var _ = Service("myservice", func() {
    Method("list", func() {
        // Requires the caller to have "read" action on "*" resource
        tscap.Require("example.com/cap/myapp", "read", "*")
        HTTP(func() { GET("/items") })
    })

    Method("create", func() {
        // Requires the caller to have "write" action on "items/*" resource
        tscap.Require("example.com/cap/myapp", "write", "items/*")
        HTTP(func() { POST("/items") })
    })

    Method("health", func() {
        // No capability check required
        tscap.AllowAnonymous()
        HTTP(func() { GET("/health") })
    })
})
```

## Matching Semantics

Grants in Tailscale ACLs can use wildcards (`*`). The DSL specifies exact requirements:

| Grant Action | Required Action | Match? |
|--------------|-----------------|--------|
| `["*"]` | `"read"` | Yes |
| `["read"]` | `"read"` | Yes |
| `["write"]` | `"read"` | No |

| Grant Resource | Required Resource | Match? |
|----------------|-------------------|--------|
| `["*"]` | `"items/123"` | Yes |
| `["items/*"]` | `"items/*"` | Yes (exact) |
| `["items/123"]` | `"items/456"` | No |

## References

- [Application capabilities](https://tailscale.com/docs/features/access-control/grants/grants-app-capabilities)
