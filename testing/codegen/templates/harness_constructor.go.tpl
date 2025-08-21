{{ printf "NewHarness creates a new test harness for the %s service." .Name | comment }}
// It starts test servers and creates a test client with transport-aware methods.
func NewHarness(t *testing.T, service {{ .PkgName }}.Service, opts ...Option) *Harness {
	t.Helper()
	
	options := &Options{
		Context: context.Background(),
		Timeout: 30 * time.Second,
	}
	
	for _, opt := range opts {
		opt(options)
	}
	
	h := &Harness{
		t:       t,
		service: service,
		ctx:     options.Context,
	}
	
	
{{- if .HasHTTP }}
	// Setup HTTP test server
	h.setupHTTP()
{{- end }}
{{- if .HasGRPC }}
	// Setup gRPC test server
	h.setupGRPC()
{{- end }}
{{- if .HasJSONRPC }}
	// Setup JSON-RPC test server
	h.setupJSONRPC()
{{- end }}
	
	// Create the test client
	h.Client = &Client{
		t:         t,
		transport: AutoTransport,
	{{- if .HasHTTP }}
		httpClient: h.getHTTPClientImpl(),
	{{- end }}
	{{- if .HasGRPC }}
		grpcClient: h.getGRPCClientImpl(),
	{{- end }}
	{{- if .HasJSONRPC }}
		jsonrpcClient: h.getJSONRPCClientImpl(),
	{{- end }}
	}
	
	// Cleanup on test completion
	t.Cleanup(func() {
		h.Close()
	})
	
	return h
}

{{ printf "Close shuts down all test servers and connections." | comment }}
func (h *Harness) Close() {
{{- if .HasHTTP }}
	if h.httpSvr != nil {
		h.httpSvr.Close()
	}
{{- end }}
{{- if .HasGRPC }}
	if h.grpcConn != nil {
		h.grpcConn.Close()
	}
	if h.grpcSvr != nil {
		h.grpcSvr.Stop()
	}
	if h.grpcLis != nil {
		h.grpcLis.Close()
	}
{{- end }}
{{- if .HasJSONRPC }}
	if h.jsonrpcSvr != nil {
		h.jsonrpcSvr.Close()
	}
{{- end }}
}
