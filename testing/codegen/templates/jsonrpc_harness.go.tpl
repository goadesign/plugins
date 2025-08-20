{{ printf "setupJSONRPC initializes the JSON-RPC test server and client." | comment }}
func (h *Harness) setupJSONRPC() {
	// Create endpoints
	endpoints := {{ .PkgName }}.NewEndpoints(h.service)
	
	// Create JSON-RPC handler
	mux := goahttp.NewMuxer()
	server := jsonrpcsvr.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil)
	jsonrpcsvr.Mount(mux, server)
	
	// Create test server
	h.jsonrpcSvr = httptest.NewServer(mux)
	
	// Create HTTP client for JSON-RPC
	h.jsonrpcCli = &http.Client{
		Timeout: 10 * time.Second,
	}
}

{{ printf "JSONRPCClient returns the HTTP client for JSON-RPC." | comment }}
func (h *Harness) JSONRPCClient() *http.Client {
	if h.jsonrpcCli == nil {
		h.t.Fatal("JSON-RPC transport not configured")
	}
	return h.jsonrpcCli
}

{{ printf "JSONRPCURL returns the URL of the JSON-RPC test server." | comment }}
func (h *Harness) JSONRPCURL() string {
	if h.jsonrpcSvr == nil {
		h.t.Fatal("JSON-RPC transport not configured")
	}
	return h.jsonrpcSvr.URL
}

{{ printf "getJSONRPCClientImpl returns the underlying JSON-RPC client implementation." | comment }}
func (h *Harness) getJSONRPCClientImpl() *jsonrpccli.Client {
	if h.jsonrpcSvr == nil || h.jsonrpcCli == nil {
		h.t.Fatal("JSON-RPC transport not configured")
	}
	u, err := url.Parse(h.jsonrpcSvr.URL)
	if err != nil {
		h.t.Fatalf("failed to parse JSON-RPC server URL: %v", err)
	}
	return jsonrpccli.NewClient(
		u.Scheme,
		u.Host,
		h.jsonrpcCli,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		false, // no debug
	)
}
