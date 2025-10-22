# Simple Tools Plugin Example

This example demonstrates how to declare a tool using the Goa tools plugin and
inspect the generated registry.

## Generate Code

```bash
go run goa.design/goa/v3/cmd/goa gen github.com/example/tools-simple/design
```

After running `goa gen` the plugin writes files under `gen/tools/`:

- `spec.go` defines the shared `ToolSpec` / `TypeSpec` structs and the
  `ToolRegistry` slice describing each defined tool.
- `codecs.go` contains the generated JSON codecs and JSON Schema blobs for every
  tool payload and result type.
