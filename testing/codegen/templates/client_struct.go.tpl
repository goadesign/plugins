{{ printf "Transport defines the transport to use for client calls." | comment }}
type Transport int

const (
	{{ printf "AutoTransport uses the first available transport (default)." | comment }}
	AutoTransport Transport = iota
	{{- if .HasHTTP }}
	{{ printf "HTTPTransport forces HTTP transport usage." | comment }}
	HTTPTransport
	{{- end }}
	{{- if .HasGRPC }}
	{{ printf "GRPCTransport forces gRPC transport usage." | comment }}
	GRPCTransport
	{{- end }}
	{{- if .HasJSONRPC }}
	{{ printf "JSONRPCTransport forces JSON-RPC transport usage." | comment }}
	JSONRPCTransport
	{{- end }}
)

{{ printf "Client provides a fluent API for calling %s service methods through various transports." .Name | comment }}
type Client struct {
	{{- if .HasHTTP }}
	httpClient *httpcli.Client
	{{- end }}
	{{- if .HasGRPC }}
	grpcClient *grpccli.Client
	{{- end }}
	{{- if .HasJSONRPC }}
	jsonrpcClient *jsonrpccli.Client
	{{- end }}
	
	// Configuration
	transport Transport
	streaming bool  // For mixed results, select streaming variant
	timeout   time.Duration
	
	// Test support
	t *testing.T
}

