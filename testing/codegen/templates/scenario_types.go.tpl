{{ printf "Scenario defines a test scenario that can be loaded from YAML." | comment }}
type Scenario struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`
	Transport   string `yaml:"transport,omitempty"` // Default transport for all steps
	Timeout     string `yaml:"timeout,omitempty"`   // Default timeout for all steps
	Steps       []Step `yaml:"steps"`
}

{{ printf "Step defines a single step in a test scenario." | comment }}
type Step struct {
	Method      string                 `yaml:"method"`
	Transport   string                 `yaml:"transport,omitempty"`   // Override transport for this step
	Payload     map[string]any `yaml:"payload,omitempty"`
	Stream      bool                   `yaml:"stream,omitempty"`      // For mixed results, use streaming variant
	Send        []map[string]any `yaml:"send,omitempty"`      // For client/bidi streaming
	Receive     []map[string]any `yaml:"receive,omitempty"`   // For server/bidi streaming
	Expect      Expectation            `yaml:"expect,omitempty"`
	Timeout     string                 `yaml:"timeout,omitempty"`     // Timeout duration (e.g., '10s', '1m')
}

{{ printf "Expectation defines expected outcomes for a step." | comment }}
type Expectation struct {
	Result map[string]any `yaml:"result,omitempty"`
	Error  string                 `yaml:"error,omitempty"`
	Stream []map[string]any `yaml:"stream,omitempty"` // Expected stream messages
	Validator string               `yaml:"validator,omitempty"` // Custom validator function name
	ValidatorPkg string            `yaml:"validator_pkg,omitempty"` // Package containing validator (defaults to service package)
}

{{ printf "ScenarioConfig holds scenarios loaded from YAML." | comment }}
type ScenarioConfig struct {
	Scenarios []Scenario `yaml:"scenarios"`
	Validators Validators `yaml:"validators,omitempty"` // Global validator configuration
}

{{ printf "Validators defines custom validation functions to use." | comment }}
type Validators struct {
	Package string `yaml:"package,omitempty"` // Default package for validators
	Path    string `yaml:"path,omitempty"`    // Import path for validator package
}

{{ printf "Valid transport values for scenarios based on service configuration." | comment }}
{{ printf "Methods may support different combinations of these transports." | comment }}
var ValidTransports = []string{
	"auto",     // Use default/first available
	{{- if .HasHTTP }}
	"http",     // HTTP plain (non-streaming methods only)
	"http-sse", // HTTP Server-Sent Events (server streaming)
	"http-ws",  // HTTP WebSocket (client/server/bidi streaming)
	{{- end }}
	{{- if .HasGRPC }}
	"grpc",     // gRPC (all streaming modes)
	{{- end }}
	{{- if .HasJSONRPC }}
	"jsonrpc",     // JSON-RPC over HTTP (non-streaming)
	"jsonrpc-sse", // JSON-RPC over SSE (server streaming)
	"jsonrpc-ws",  // JSON-RPC over WebSocket (streaming only)
	{{- end }}
}

{{ printf "TransportAvailability documents which transports each method supports." | comment }}
var TransportAvailability = map[string][]string{
	{{- range .Methods }}
	"{{ .Name }}": { {{- range $i, $t := .Transports }}{{- if $i }}, {{ end }}"{{ $t }}"{{- end }} },
	{{- end }}
}