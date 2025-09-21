// Package mains defines a Goa plugin that customizes the example server
// main layout to follow the recommended conventions.
//
// Behavior:
//   - Single-service servers: writes services/<svc>/cmd/<svc>/main.go
//   - Multi-service servers:  writes cmd/<server>/main.go
//   - Removes default cmd/<server>/main.go and cmd/<server>/http.go
//   - Wires Clue + OpenTelemetry, debug endpoints, and health/metrics server
//   - Detects WebSocket use (non-SSE streaming) and includes an upgrader only
//     when needed
//
// Usage:
//   // Import the plugin to register it at init time.
//   import _ "goa.design/plugins/v3/mains"
//
// Then run:
//   goa gen     <module>/design
//   goa example <module>/design
package mains

