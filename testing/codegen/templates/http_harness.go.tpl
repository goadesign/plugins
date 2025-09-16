{{ printf "setupHTTP initializes the HTTP test server and client." | comment }}
func (h *Harness) setupHTTP() {
	// Create endpoints
	endpoints := {{ .PkgName }}.NewEndpoints(h.service)
	
	// Create HTTP handler
	mux := goahttp.NewMuxer()
	{{- if .HasStreams }}
	// Create WebSocket upgrader for streaming endpoints
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	{{- end }}
	server := httpsvr.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil{{ if .HasStreams }}, upgrader, nil{{ end }})
	httpsvr.Mount(mux, server)
	
	// Create test server
	h.httpSvr = httptest.NewServer(mux)
	
	// Create HTTP client
	h.httpCli = &http.Client{
		Timeout: 10 * time.Second,
	}
}

{{ printf "HTTPClient returns an HTTP client configured for the test server." | comment }}
func (h *Harness) HTTPClient() *http.Client {
	if h.httpCli == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	return h.httpCli
}

{{ printf "getHTTPClientImpl returns the underlying HTTP client implementation." | comment }}
func (h *Harness) getHTTPClientImpl() *httpcli.Client {
	if h.httpSvr == nil || h.httpCli == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	u, err := url.Parse(h.httpSvr.URL)
	if err != nil {
		h.t.Fatalf("invalid test server URL: %v", err)
	}
	scheme := u.Scheme
	host := u.Host
	{{- if .HasStreams }}
	// Create WebSocket dialer for streaming endpoints
	wsDialer := &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
	}
	{{- end }}

	return httpcli.NewClient(
		scheme,
		host,
		h.httpCli,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		false,
		{{- if .HasStreams }}
		wsDialer,
		nil,
		{{- end }}
	)
}

{{ printf "HTTPClientEndpoints creates HTTP client endpoints for the service." | comment }}
func (h *Harness) HTTPClientEndpoints() *{{ .PkgName }}.Endpoints {
	c := h.getHTTPClientImpl()
	return &{{ .PkgName }}.Endpoints{
		{{- range .Methods }}
		{{- $method := . }}
		{{- range .Targets }}
		{{- if or .IsHTTPPlain .IsHTTPServerSent .IsHTTPWebSocket }}
		{{ $method.VarName }}: c.{{ $method.VarName }}(),
		{{- break }}
		{{- end }}
		{{- end }}
		{{- end }}
	}
}

{{ printf "HTTPURL returns the base URL of the test HTTP server." | comment }}
func (h *Harness) HTTPURL() string {
	if h.httpSvr == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	return h.httpSvr.URL
}

{{ printf "HTTPRequest creates a new HTTP request for testing." | comment }}
func (h *Harness) HTTPRequest(method, path string, body any) *http.Request {
	h.t.Helper()
	
	url := h.HTTPURL() + path
	
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			h.t.Fatalf("Failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	
	req, err := http.NewRequestWithContext(h.ctx, method, url, bodyReader)
	if err != nil {
		h.t.Fatalf("Failed to create request: %v", err)
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req
}

{{ printf "HTTPWSURL builds a websocket URL from the base HTTP URL and path." | comment }}
func (h *Harness) HTTPWSURL(path string) string {
	base := h.HTTPURL()
	u, err := url.Parse(base)
	if err != nil {
		h.t.Fatalf("invalid base URL: %v", err)
	}
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}
	u.Path = path
	return u.String()
}

{{ printf "HTTPDo performs an HTTP request and returns the response." | comment }}
func (h *Harness) HTTPDo(req *http.Request) *http.Response {
	h.t.Helper()
	
	resp, err := h.HTTPClient().Do(req)
	if err != nil {
		h.t.Fatalf("HTTP request failed: %v", err)
	}
	
	return resp
}
