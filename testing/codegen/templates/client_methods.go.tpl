{{- range .Methods }}
{{- $method := . }}

{{ printf "%s calls the %s method using the configured transport." $method.VarName $method.Name | comment }}
{{- if or $method.ServerStream $method.ClientStream }}
func (c *Client) {{ $method.VarName }}(ctx context.Context{{- if and $method.PayloadRef (not $method.StreamingPayload) }}, p *{{ $.PkgName }}.{{ $method.Payload }}{{- end }}) ({{ $.PkgName }}.{{ $method.ClientStream.Interface }}, error) {
{{- else }}
func (c *Client) {{ $method.VarName }}(ctx context.Context{{- if $method.PayloadRef }}, p *{{ $.PkgName }}.{{ $method.Payload }}{{- end }}) ({{- if $method.PkgResultRef }}{{ $method.PkgResultRef }}, {{ end }}error) {
{{- end }}
	// Determine which transport to use
	transport := c.transport
	if transport == AutoTransport {
		// Use the first available transport
		{{- if .HasHTTP }}
		if c.httpClient != nil {
			transport = HTTPTransport
		}{{- if .HasGRPC }} else{{- end }}
		{{- end }}
		{{- if .HasGRPC }}
		if c.grpcClient != nil {
			transport = GRPCTransport
		}{{- if .HasJSONRPC }} else{{- end }}
		{{- end }}
		{{- if .HasJSONRPC }}
		if c.jsonrpcClient != nil {
			transport = JSONRPCTransport
		}
		{{- end }}
	}
	
	switch transport {
	{{- if .HasHTTP }}
	case HTTPTransport:
		if c.httpClient == nil {
			return {{ if or $method.ResultRef $method.ClientStream }}nil, {{ end }}fmt.Errorf("HTTP transport not configured")
		}
		endpoint := c.httpClient.{{ $method.VarName }}()
		{{- if or $method.ServerStream $method.ClientStream }}
		res, err := endpoint(ctx{{- if and $method.PayloadRef (not $method.StreamingPayload) }}, p{{ else }}, nil{{ end }})
		if err != nil {
			return nil, err
		}
		return res.({{ $.PkgName }}.{{ $method.ClientStream.Interface }}), nil
		{{- else }}
		{{- if $method.ResultRef }}
		res, err := endpoint(ctx{{- if $method.PayloadRef }}, p{{ end }})
		if err != nil {
			return nil, err
		}
		return res.({{ $method.PkgResultRef }}), nil
		{{- else }}
		_, err := endpoint(ctx, {{- if $method.PayloadRef }}p{{ else }}nil{{ end }})
		return err
		{{- end }}
		{{- end }}
	{{- end }}
	
	{{- if .HasGRPC }}
	case GRPCTransport:
		if c.grpcClient == nil {
			return {{ if or $method.ResultRef $method.ClientStream }}nil, {{ end }}fmt.Errorf("gRPC transport not configured")
		}
		endpoint := c.grpcClient.{{ $method.VarName }}()
		{{- if or $method.ServerStream $method.ClientStream }}
		res, err := endpoint(ctx{{- if and $method.PayloadRef (not $method.StreamingPayload) }}, p{{ else }}, nil{{ end }})
		if err != nil {
			return nil, err
		}
		return res.({{ $.PkgName }}.{{ $method.ClientStream.Interface }}), nil
		{{- else }}
		{{- if $method.ResultRef }}
		res, err := endpoint(ctx{{- if $method.PayloadRef }}, p{{ end }})
		if err != nil {
			return nil, err
		}
		return res.({{ $method.PkgResultRef }}), nil
		{{- else }}
		_, err := endpoint(ctx, {{- if $method.PayloadRef }}p{{ else }}nil{{ end }})
		return err
		{{- end }}
		{{- end }}
	{{- end }}
	
	{{- if .HasJSONRPC }}
	case JSONRPCTransport:
		if c.jsonrpcClient == nil {
			return {{ if or $method.ResultRef $method.ClientStream }}nil, {{ end }}fmt.Errorf("JSON-RPC transport not configured")
		}
		endpoint := c.jsonrpcClient.{{ $method.VarName }}()
		{{- if or $method.ServerStream $method.ClientStream }}
		res, err := endpoint(ctx{{- if and $method.PayloadRef (not $method.StreamingPayload) }}, p{{ else }}, nil{{ end }})
		if err != nil {
			return nil, err
		}
		{{- /* Check if this is a JSON-RPC SSE stream that needs wrapping */ -}}
		{{- $isJSONRPCSSE := false }}
		{{- range .Targets }}
			{{- if .IsJSONRPCSSE }}
				{{- $isJSONRPCSSE = true }}
			{{- end }}
		{{- end }}
		{{- if $isJSONRPCSSE }}
		// Wrap JSON-RPC SSE stream to match service interface
		return &jsonrpcSSE{{ $method.VarName }}Wrapper{stream: res.(interface {
			Recv(context.Context) (*{{ $.PkgName }}.{{ $method.Result }}, error)
			Close() error
		})}, nil
		{{- else }}
		return res.({{ $.PkgName }}.{{ $method.ClientStream.Interface }}), nil
		{{- end }}
		{{- else }}
		{{- if $method.ResultRef }}
		res, err := endpoint(ctx{{- if $method.PayloadRef }}, p{{ end }})
		if err != nil {
			return nil, err
		}
		return res.({{ $method.PkgResultRef }}), nil
		{{- else }}
		_, err := endpoint(ctx, {{- if $method.PayloadRef }}p{{ else }}nil{{ end }})
		return err
		{{- end }}
		{{- end }}
	{{- end }}
	
	default:
		return {{ if or $method.ResultRef $method.ClientStream }}nil, {{ end }}fmt.Errorf("no transport available for {{ $method.Name }}")
	}
}

{{- end }}
