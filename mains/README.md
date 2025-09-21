# mains plugin

The `mains` plugin customizes the `goa example` generated server mains to follow the Goa golden-path layout and wiring.

Features:
- Single-service servers: writes `services/<svc>/cmd/<svc>/main.go`
- Multi-service servers: writes `cmd/<server>/main.go`
- Clue + OpenTelemetry wiring, debug endpoints, health + metrics server
- WebSocket-aware: includes `websocket.Upgrader` only when needed

Usage:

```go
import _ "goa.design/plugins/v3/mains" // register plugin
```

Then run:

```
goa gen <module>/design
goa example <module>/design
```

The plugin inspects the example-generated HTTP server files and the design roots; no DSL import or extra configuration is required.

