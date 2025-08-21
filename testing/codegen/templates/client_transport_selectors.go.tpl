{{- if .HasHTTP }}
{{ printf "HTTP returns a new client configured to use HTTP transport." | comment }}
func (c *Client) HTTP() *Client {
	nc := *c
	nc.transport = HTTPTransport
	return &nc
}
{{- end }}

{{- if .HasGRPC }}
{{ printf "GRPC returns a new client configured to use gRPC transport." | comment }}
func (c *Client) GRPC() *Client {
	nc := *c
	nc.transport = GRPCTransport
	return &nc
}
{{- end }}

{{- if .HasJSONRPC }}
{{ printf "JSONRPC returns a new client configured to use JSON-RPC transport." | comment }}
func (c *Client) JSONRPC() *Client {
	nc := *c
	nc.transport = JSONRPCTransport
	return &nc
}
{{- end }}

{{ printf "AsStream returns a new client configured to request streaming variants for mixed results." | comment }}
{{ printf "This is only relevant for methods with SSE endpoints that support content negotiation." | comment }}
func (c *Client) AsStream() *Client {
	nc := *c
	nc.streaming = true
	return &nc
}

{{ printf "WithTimeout returns a new client with the specified timeout." | comment }}
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	nc := *c
	nc.timeout = timeout
	return &nc
}