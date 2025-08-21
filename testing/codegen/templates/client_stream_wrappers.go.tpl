{{- /* Stream wrapper types for adapting transport-specific streams to service interfaces */ -}}

{{- range .Methods }}
{{- $method := . }}
{{- if $method.ServerStream }}
{{- range .Targets }}
{{- if .IsHTTPServerSent }}
{{ printf "httpSSE%sWrapper wraps HTTP SSE stream to match service interface." $method.VarName | comment }}
type httpSSE{{ $method.VarName }}Wrapper struct {
	stream interface {
		Recv(context.Context) (*{{ $.PkgName }}.{{ $method.Result }}, error)
		Close() error
	}
}

// Recv implements the service interface.
func (w *httpSSE{{ $method.VarName }}Wrapper) Recv() (*{{ $.PkgName }}.{{ $method.Result }}, error) {
	return w.stream.Recv(context.Background())
}

// RecvWithContext implements the service interface.
func (w *httpSSE{{ $method.VarName }}Wrapper) RecvWithContext(ctx context.Context) (*{{ $.PkgName }}.{{ $method.Result }}, error) {
	return w.stream.Recv(ctx)
}
{{- end }}
{{- if .IsJSONRPCSSE }}
{{ printf "jsonrpcSSE%sWrapper wraps JSON-RPC SSE stream to match service interface." $method.VarName | comment }}
type jsonrpcSSE{{ $method.VarName }}Wrapper struct {
	stream interface {
		Recv(context.Context) (*{{ $.PkgName }}.{{ $method.Result }}, error)
		Close() error
	}
}

// Recv implements the service interface.
func (w *jsonrpcSSE{{ $method.VarName }}Wrapper) Recv() (*{{ $.PkgName }}.{{ $method.Result }}, error) {
	return w.stream.Recv(context.Background())
}

// RecvWithContext implements the service interface.
func (w *jsonrpcSSE{{ $method.VarName }}Wrapper) RecvWithContext(ctx context.Context) (*{{ $.PkgName }}.{{ $method.Result }}, error) {
	return w.stream.Recv(ctx)
}
{{- end }}
{{- end }}
{{- end }}
{{- end }}