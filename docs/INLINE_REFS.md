## Goal

Design a robust, deterministic strategy to inline all JSON Schema $ref occurrences into their target schemas for docs JSON produced by goa. The generator only uses local references to the top-level definitions object (paths of the form "#/definitions/<Name>").

## Where $ref appear in goa’s output

- **User/Result types**: Type schemas reference definitions via `TypeRef*` and `ResultTypeRef*`.
- **Object properties**: Nested fields can be `$ref`.
- **Array items**: `items` may be a `$ref` when the element type is a user/result type.
- **Map values**: `additionalProperties` may be a schema containing `$ref` (or boolean `true`).
- **Union types**: `anyOf` is an array of schemas that may contain `$ref`.
- **Service-level top schema**: APISchema properties refer to `#/definitions/<Service>`.

Assumption: Only references to `#/definitions/<Name>` are present. No external files/URLs.

## Constraints

- Definitions are in a per-root map; keys are type names.
- Definitions can themselves contain nested `$ref`.
- Must avoid infinite recursion on cyclical definitions.
- Preserve required fields, examples, and order where applicable.

## Strategy (high-level)

1) **Snapshot definitions**
   - Work on a deep-copied map of the per-root `definitions` to avoid mutating shared state.

2) **Inlining pass**
   - Implement `inlineRefs(schema, defs, stack)`:
     - If `schema.Ref == "#/definitions/<Name>"`:
       - If `<Name>` is in `stack`, keep `$ref` (cycle break) and return.
       - Else: push `<Name>`, deep-copy `defs[Name]`, recursively `inlineRefs` on the copy, then replace current node with the copy and clear `Ref`. Pop `<Name>`.
     - If `schema.Ref == ""`: recursively process all composite positions:
       - `properties` values
       - `items`
       - `additionalProperties` when it is a `*Schema`
       - `anyOf` elements
       - (If present) nested `definitions` for completeness

3) **Apply order**
   - Perform schema transforms that affect keys first (e.g., JSON tag rename and required filtering), then inline.

4) **Service-level application**
   - Run `inlineRefs` on every reachable schema under services (payload, result, streaming payload/result, error types).
   - Optionally inline inside `definitions` if you plan to drop them entirely.

5) **Cycles**
   - Use a `stack` (set of definition names being expanded). If a name is already in the stack, retain `$ref` at that edge to prevent infinite expansion. This yields minimal `$ref` for strongly-cyclic graphs.

6) **AnyOf**
   - Goa builds `anyOf` by appending union variants in-order. Inline each element; preserve order.

7) **Maps**
   - If `additionalProperties` is `true`, leave as is. If it is a schema, inline there as well.

8) **Performance**
   - Start without memoization. If profiling reveals hotspots, consider caching fully inlined, deep-copied definitions by name and reusing them, taking care not to share mutable pointers unexpectedly.

9) **Post-processing**
   - If you require a completely `$ref`-free document, remove `definitions` after inlining all reachable schemas. Otherwise, keep or prune unused definitions based on tests/goldens.

## Pseudocode

```go
func inlineAllServiceSchemas(d *data, defs map[string]*openapi.Schema) {
    stack := map[string]bool{}
    for _, s := range d.Services {
        for _, m := range s.Methods {
            if m.Payload != nil && m.Payload.Type != nil {
                inlineRefs(m.Payload.Type, defs, stack)
            }
            if m.StreamingPayload != nil && m.StreamingPayload.Type != nil {
                inlineRefs(m.StreamingPayload.Type, defs, stack)
            }
            if m.Result != nil && m.Result.Type != nil {
                inlineRefs(m.Result.Type, defs, stack)
            }
            if m.StreamingResult != nil && m.StreamingResult.Type != nil {
                inlineRefs(m.StreamingResult.Type, defs, stack)
            }
            for _, e := range m.Errors {
                if e.Type != nil {
                    inlineRefs(e.Type, defs, stack)
                }
            }
        }
    }
}

func inlineRefs(s *openapi.Schema, defs map[string]*openapi.Schema, stack map[string]bool) {
    if s == nil {
        return
    }
    if s.Ref != "" {
        const prefix = "#/definitions/"
        if !strings.HasPrefix(s.Ref, prefix) { return }
        name := strings.TrimPrefix(s.Ref, prefix)
        if stack[name] { return }
        def, ok := defs[name]
        if !ok || def == nil { return }
        stack[name] = true
        copy := dupSchema(def)
        inlineRefs(copy, defs, stack)
        *s = *copy
        s.Ref = ""
        delete(stack, name)
        return
    }
    for _, p := range s.Properties { inlineRefs(p, defs, stack) }
    if s.Items != nil { inlineRefs(s.Items, defs, stack) }
    if ap, ok := s.AdditionalProperties.(*openapi.Schema); ok && ap != nil {
        inlineRefs(ap, defs, stack)
    }
    for _, a := range s.AnyOf { inlineRefs(a, defs, stack) }
    for _, d := range s.Definitions { inlineRefs(d, defs, stack) }
}
```

## Integration points in the plugin

- After building `docs` and `defs`, and after JSON tag transforms:
  - Gate behind `InlineRefs` option: call `inlineAllServiceSchemas(docs, defs)`.
  - Decide whether to keep or drop `docs.Definitions` based on goldens. If keeping, you may optionally prune unused definitions.

## Edge cases and guarantees

- **Cycles**: Minimal `$ref` retained on cycle entry.
- **Required and examples**: Preserved; JSON-tag transform already remaps them pre-inlining.
- **Maps**: Free-form (`true`) untouched; schema-valued inlined.
- **Unions**: Order preserved in `anyOf`.

## Rationale

- Tailored to the shapes produced by goa: only `#/definitions/` refs.
- Deterministic and safe (deep copies, cycle guard), compositional (handles all container fields).
- Clean pipeline: rename/tag first, inline second.


