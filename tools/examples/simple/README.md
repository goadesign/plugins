# Simple Tools Plugin Example

This example demonstrates how to declare a tool using the Goa tools plugin and
inspect the generated registry.

## Generate Code

```bash
go run goa.design/goa/v3/cmd/goa gen github.com/example/tools-simple/design
```

After running `goa gen` the plugin writes service-scoped tool packages, for
example:

- `gen/inventory/tools/inventory_tools/` for the pure tools declared in the
  design. This package contains `types.go`, `schemas.go`, `codecs.go`, and
  `registry.go` for the `lookup_item` and `list_recent_items` tools.
- `gen/inventory/tools/inventory_method_tools/` for the method-derived tools.
  Only codecs, schemas, and the registry are generated because the payload and
  result structs live in the service package already.
