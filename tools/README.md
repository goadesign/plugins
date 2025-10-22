# Goa Tools Plugin

Many systems built with Goa expose **tools**: structured commands that other
components (LLM workflows, background jobs, admin dashboards) invoke. A tool has
a name plus validated payload and result schemas. Without a central definition it
is easy for services to drift—each consumer recreates its own structs, codecs,
and validation rules.

The Tools plugin keeps those contracts in sync. Add the DSL to your Goa design,
and the plugin generates a single, canonical set of artefacts—Go structs, JSON
Schema, JSON codecs, and a registry—that every consumer can share.

## Why use it?

Use the plugin whenever the same tool must flow through multiple subsystems:

- LLM orchestration or chat loops that call back into services.
- Temporal payload converters, session persistence, or job workers.
- Administrative or diagnostics tooling that inspects past tool invocations.
- Multiple services that would otherwise hand-roll serializers.

Define the tool once, regenerate, and every consumer speaks the same schema.

## Generated assets

Running `goa gen` produces, for each tool set, a dedicated package under
`gen/<service>/tools/<toolset>/` (or `gen/tools/<toolset>/` when the set isn't
service-scoped). Each package contains:

| File          | Purpose                                                                                           |
|---------------|---------------------------------------------------------------------------------------------------|
| `types.go`    | Go structs for pure tools defined in the set. Method-derived tools reuse the service codegen types. |
| `schemas.go`  | Embedded JSON Schema literals for every payload/result.                                           |
| `codecs.go`   | Typed marshal/unmarshal helpers with inlined Goa validation, plus name-driven generic codecs.     |
| `registry.go` | A catalogue of `tools.ToolSpec` entries (name, service, set, codecs, schemas) with helper queries. |

Runtime helpers (`tools/codec.go`, `tools/toolspec.go`) live in this module, so
generated code simply imports `goa.design/plugins/v3/tools` alongside the
tool-set package that was generated for the service.

## DSL overview

```go
import . "goa.design/plugins/v3/tools/dsl"

Service("inventory", func() {
    ToolSet("ops", func() {
        Tool("lookup_item", func() {
            Payload(func() {
                Attribute("sku", String, "Inventory SKU", func() { MinLength(1) })
                Required("sku")
            })
            Result(func() {
                Attribute("found", Boolean, "True when item exists")
                Attribute("description", String, "Optional description")
                Required("found")
            })
        })

        ToolFromMethod("ListRecentItems") // reuse existing Goa method
    })
})
```

- `ToolSet(name, func())` scopes tools to a Goa service.
- `Tool(name, func())` defines a standalone tool (the DSL mirrors Goa `Method`).
- `ToolFromMethod(method, optionalName)` reuses an existing service method so
you only maintain one payload/result definition.

Pure tool bodies execute during DSL evaluation so the generator sees the final
payload/result attributes. Method-derived tools are tagged so the collector can
reuse the service-generated types, including custom `struct:pkg:path` locations.

## Consuming the generated artefacts

```go
import inventorytools "github.com/example/project/gen/inventory/tools/inventory_tools"

codec, ok := inventorytools.PayloadCodec("lookup_item")
if !ok {
    return fmt.Errorf("unknown tool")
}
raw, err := codec.ToJSON(&inventorytools.LookupItemPayload{Sku: "ABC"})
```

- `PayloadCodec` / `ResultCodec` provide generic JSON codecs backed by the typed helpers.
- Typed helpers (`MarshalLookupItemPayload`, `UnmarshalLookupItemPayload`, `ValidateLookupItemPayload`, …)
  include the Goa-generated validation logic.
- `ToolRegistry`, `Names`, `Spec`, `PayloadSchema`, `ResultSchema` expose
  metadata for admin or debugging views.

## Development & tests

- `go test ./tools/...` exercises the generator and DSL integration.
- `tools/examples/simple` demonstrates a minimal design; run
  `goa gen github.com/example/tools-simple/design` in that directory to inspect
  the output.

## FAQ

**What is a tool?**  A named command with a validated payload/result that other
components can call (e.g. `lookup_item`).

**Why not rely on Goa methods alone?**  Goa methods generate transport handlers
(HTTP/gRPC). Tools are often transport-agnostic—invoked by LLMs or background
workers—and may not have direct endpoints. The plugin lets you model those
contracts without adding extra handlers.

**Do I always get new structs?**  Pure tools (defined via `Tool`) generate new
structs. Tools created with `ToolFromMethod` reuse the service’s existing
payload/result types, including custom `struct:pkg:path` locations.
