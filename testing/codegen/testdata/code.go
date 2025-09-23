package testdata

var WithResultCode = `// WithResultMethod calls the WithResultMethod method using the configured
// transport.
func (c *Client) WithResultMethod(ctx context.Context) (*withresultservice.WithResultMethodResult, error) {
	// Determine which transport to use
	transport := c.transport
	if transport == AutoTransport {
		// Use the first available transport
		if c.httpClient != nil {
			transport = HTTPTransport
		} else if c.grpcClient != nil {
			transport = GRPCTransport
		} else if c.jsonrpcClient != nil {
			transport = JSONRPCTransport
		}
	}

	switch transport {
	case HTTPTransport:
		if c.httpClient == nil {
			return nil, fmt.Errorf("HTTP transport not configured")
		}
		endpoint := c.httpClient.WithResultMethod()
		res, err := endpoint(ctx)
		if err != nil {
			return nil, err
		}
		return res.(*withresultservice.WithResultMethodResult), nil
	case GRPCTransport:
		if c.grpcClient == nil {
			return nil, fmt.Errorf("gRPC transport not configured")
		}
		endpoint := c.grpcClient.WithResultMethod()
		res, err := endpoint(ctx)
		if err != nil {
			return nil, err
		}
		return res.(*withresultservice.WithResultMethodResult), nil
	case JSONRPCTransport:
		if c.jsonrpcClient == nil {
			return nil, fmt.Errorf("JSON-RPC transport not configured")
		}
		endpoint := c.jsonrpcClient.WithResultMethod()
		res, err := endpoint(ctx)
		if err != nil {
			return nil, err
		}
		return res.(*withresultservice.WithResultMethodResult), nil

	default:
		return nil, fmt.Errorf("no transport available for WithResultMethod")
	}
}
`

var WithoutResultCode = `// WithoutResultMethod calls the WithoutResultMethod method using the
// configured transport.
func (c *Client) WithoutResultMethod(ctx context.Context) error {
	// Determine which transport to use
	transport := c.transport
	if transport == AutoTransport {
		// Use the first available transport
		if c.httpClient != nil {
			transport = HTTPTransport
		} else if c.grpcClient != nil {
			transport = GRPCTransport
		} else if c.jsonrpcClient != nil {
			transport = JSONRPCTransport
		}
	}

	switch transport {
	case HTTPTransport:
		if c.httpClient == nil {
			return fmt.Errorf("HTTP transport not configured")
		}
		endpoint := c.httpClient.WithoutResultMethod()
		return endpoint(ctx)
	case GRPCTransport:
		if c.grpcClient == nil {
			return fmt.Errorf("gRPC transport not configured")
		}
		endpoint := c.grpcClient.WithoutResultMethod()
		return endpoint(ctx)
	case JSONRPCTransport:
		if c.jsonrpcClient == nil {
			return fmt.Errorf("JSON-RPC transport not configured")
		}
		endpoint := c.jsonrpcClient.WithoutResultMethod()
		return endpoint(ctx)

	default:
		return fmt.Errorf("no transport available for WithoutResultMethod")
	}
}
`

var WithStreamCode = `// setupHTTP initializes the HTTP test server and client.
func (h *Harness) setupHTTP() {
	// Create endpoints
	endpoints := withstreamservice.NewEndpoints(h.service)

	// Create HTTP handler
	mux := goahttp.NewMuxer()
	// Create WebSocket upgrader for streaming endpoints
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	server := httpsvr.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil, upgrader, nil)
	httpsvr.Mount(mux, server)

	// Create test server
	h.httpSvr = httptest.NewServer(mux)

	// Create HTTP client
	h.httpCli = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// HTTPClient returns an HTTP client configured for the test server.
func (h *Harness) HTTPClient() *http.Client {
	if h.httpCli == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	return h.httpCli
}

// getHTTPClientImpl returns the underlying HTTP client implementation.
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
	// Create WebSocket dialer for streaming endpoints
	wsDialer := &websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
	}

	return httpcli.NewClient(
		scheme,
		host,
		h.httpCli,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		false,
		wsDialer,
		nil,
	)
}

// HTTPClientEndpoints creates HTTP client endpoints for the service.
func (h *Harness) HTTPClientEndpoints() *withstreamservice.Endpoints {
	c := h.getHTTPClientImpl()
	return &withstreamservice.Endpoints{
		WithStreamMethod: c.WithStreamMethod(),
	}
}

// HTTPURL returns the base URL of the test HTTP server.
func (h *Harness) HTTPURL() string {
	if h.httpSvr == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	return h.httpSvr.URL
}

// HTTPRequest creates a new HTTP request for testing.
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

// HTTPWSURL builds a websocket URL from the base HTTP URL and path.
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

// HTTPDo performs an HTTP request and returns the response.
func (h *Harness) HTTPDo(req *http.Request) *http.Response {
	h.t.Helper()

	resp, err := h.HTTPClient().Do(req)
	if err != nil {
		h.t.Fatalf("HTTP request failed: %v", err)
	}

	return resp
}
`

var WithoutStreamCode = `// setupHTTP initializes the HTTP test server and client.
func (h *Harness) setupHTTP() {
	// Create endpoints
	endpoints := withoutstreamservice.NewEndpoints(h.service)

	// Create HTTP handler
	mux := goahttp.NewMuxer()
	server := httpsvr.New(endpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	httpsvr.Mount(mux, server)

	// Create test server
	h.httpSvr = httptest.NewServer(mux)

	// Create HTTP client
	h.httpCli = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// HTTPClient returns an HTTP client configured for the test server.
func (h *Harness) HTTPClient() *http.Client {
	if h.httpCli == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	return h.httpCli
}

// getHTTPClientImpl returns the underlying HTTP client implementation.
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

	return httpcli.NewClient(
		scheme,
		host,
		h.httpCli,
		goahttp.RequestEncoder,
		goahttp.ResponseDecoder,
		false,
	)
}

// HTTPClientEndpoints creates HTTP client endpoints for the service.
func (h *Harness) HTTPClientEndpoints() *withoutstreamservice.Endpoints {
	c := h.getHTTPClientImpl()
	return &withoutstreamservice.Endpoints{
		WithoutStreamMethod: c.WithoutStreamMethod(),
	}
}

// HTTPURL returns the base URL of the test HTTP server.
func (h *Harness) HTTPURL() string {
	if h.httpSvr == nil {
		h.t.Fatal("HTTP transport not configured")
	}
	return h.httpSvr.URL
}

// HTTPRequest creates a new HTTP request for testing.
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

// HTTPDo performs an HTTP request and returns the response.
func (h *Harness) HTTPDo(req *http.Request) *http.Response {
	h.t.Helper()

	resp, err := h.HTTPClient().Do(req)
	if err != nil {
		h.t.Fatalf("HTTP request failed: %v", err)
	}

	return resp
}
`
