package testdata

var ClientMethodsWithResultCode = `// WithResultMethod calls the WithResultMethod method using the configured
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

var ClientMethodsWithoutResultCode = `// WithoutResultMethod calls the WithoutResultMethod method using the
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
		_, err := endpoint(ctx, nil)
		return err
	case GRPCTransport:
		if c.grpcClient == nil {
			return fmt.Errorf("gRPC transport not configured")
		}
		endpoint := c.grpcClient.WithoutResultMethod()
		_, err := endpoint(ctx, nil)
		return err
	case JSONRPCTransport:
		if c.jsonrpcClient == nil {
			return fmt.Errorf("JSON-RPC transport not configured")
		}
		endpoint := c.jsonrpcClient.WithoutResultMethod()
		_, err := endpoint(ctx, nil)
		return err

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

var ScenarioRunnerWithResultCode = `// ScenarioRunner executes test scenarios.
type ScenarioRunner struct {
	scenarios  []Scenario
	validators Validators // Global validator configuration
}

// LoadScenarios loads scenarios from a YAML file.
func LoadScenarios(path string) (*ScenarioRunner, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios file: %w", err)
	}

	var config ScenarioConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse scenarios YAML: %w", err)
	}

	return &ScenarioRunner{
		scenarios:  config.Scenarios,
		validators: config.Validators,
	}, nil
}

// NewScenarioRunner creates a new scenario runner.
func NewScenarioRunner() *ScenarioRunner {
	return &ScenarioRunner{
		scenarios: make([]Scenario, 0),
	}
}

// AddScenario adds a scenario to the runner.
func (r *ScenarioRunner) AddScenario(scenario Scenario) {
	r.scenarios = append(r.scenarios, scenario)
}

// Run executes all scenarios.
func (r *ScenarioRunner) Run(t *testing.T, client *Client) {
	if r == nil {
		t.Fatal("ScenarioRunner is nil")
	}
	if client == nil {
		t.Fatal("Client is nil")
	}
	for _, scenario := range r.scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			r.runScenario(t, client, scenario)
		})
	}
}

// RunNamed executes a specific scenario by name.
func (r *ScenarioRunner) RunNamed(t *testing.T, client *Client, name string) {
	if r == nil {
		t.Fatal("ScenarioRunner is nil")
	}
	if client == nil {
		t.Fatal("Client is nil")
	}
	if name == "" {
		t.Fatal("scenario name is empty")
	}
	for _, scenario := range r.scenarios {
		if scenario.Name == name {
			r.runScenario(t, client, scenario)
			return
		}
	}
	t.Fatalf("scenario %q not found", name)
}

func (r *ScenarioRunner) runScenario(t *testing.T, client *Client, scenario Scenario) {
	// Apply default transport if specified
	if scenario.Transport != "" {
		client = r.selectTransport(client, scenario.Transport)
	}

	for i, step := range scenario.Steps {
		t.Run(fmt.Sprintf("step_%d_%s", i+1, step.Method), func(t *testing.T) {
			// Apply scenario-level timeout if step doesn't override
			if step.Timeout == "" && scenario.Timeout != "" {
				step.Timeout = scenario.Timeout
			}
			r.runStep(t, client, step)
		})
	}
}

func (r *ScenarioRunner) runStep(t *testing.T, client *Client, step Step) {
	// Apply per-step transport override
	if step.Transport != "" {
		client = r.selectTransport(client, step.Transport)
	}

	// Validate transport availability
	if step.Transport != "" && step.Transport != "auto" {
		if transports, ok := TransportAvailability[step.Method]; ok {
			found := false
			for _, t := range transports {
				if t == step.Transport {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("method %q does not support transport %q, available: %v",
					step.Method, step.Transport, transports)
			}
		}
	}

	// Process payload
	payload := step.Payload
	ctx := context.Background()

	// Apply timeout if specified
	if step.Timeout != "" {
		duration, err := time.ParseDuration(step.Timeout)
		if err != nil {
			t.Fatalf("invalid timeout %q: %v", step.Timeout, err)
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, duration)
		defer cancel()
	}

	// Execute the method
	result, err := r.executeMethod(ctx, client, step.Method, payload)

	// Handle error expectation
	if step.Expect.Error != "" {
		if err == nil {
			t.Errorf("expected error %q but got none", step.Expect.Error)
		} else if !strings.Contains(err.Error(), step.Expect.Error) {
			t.Errorf("expected error containing %q but got %q", step.Expect.Error, err.Error())
		}
		return
	}

	// Handle unexpected error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Validate result if expected
	if step.Expect.Result != nil || step.Expect.Validator != "" {
		r.validateResult(t, step.Method, result, step.Expect)
	}

	// Handle streaming expectations
	if len(step.Expect.Stream) > 0 {
		r.validateStream(t, step.Method, result, step.Expect)
	}
}

func (r *ScenarioRunner) executeMethod(ctx context.Context, client *Client, method string, payload map[string]any) (any, error) {
	switch method {
	case "WithResultMethod":
		return client.WithResultMethod(ctx)
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

func (r *ScenarioRunner) mapToStruct(data map[string]any, target any) error {
	if data == nil {
		// nil data is okay, just return without setting anything
		return nil
	}
	if target == nil {
		return fmt.Errorf("target is nil")
	}
	// Convert map to JSON then unmarshal to struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}

func (r *ScenarioRunner) validateResult(t *testing.T, method string, result any, expect Expectation) {
	if result == nil && expect.Result == nil && expect.Validator == "" {
		// Nothing to validate
		return
	}

	// If custom validator specified in YAML, call it
	if expect.Validator != "" {
		// Call the user-defined validator function
		// The function signature should be: func(t *testing.T, result *ServiceType, expected map[string]any)
		r.callValidator(t, method, result, expect)
		return
	}

	// Fall back to default validation
	if expect.Result != nil {
		if err := defaultValidateResult(result, expect.Result); err != nil {
			t.Errorf("validation failed for %s: %v", method, err)
		}
	}
}

// callValidator calls the user-specified validator function.
// The validator function must be defined in the test package.
func (r *ScenarioRunner) callValidator(t *testing.T, method string, result any, expect Expectation) {
	// For each validator found in YAML, we generate a direct call
	// Users must define these functions in their test files

	validatorName := expect.Validator
	_ = validatorName // avoid unused variable in case no validators are defined

	switch method {
	case "WithResultMethod":
		typedResult := result.(*withresultservice.WithResultMethodResult)
		_ = typedResult // no validators defined in YAML
		t.Errorf("validator %q specified but not generated - add it to scenarios.yaml first", validatorName)
	default:
		t.Errorf("unknown method: %s", method)
	}
}

// defaultValidateResult provides basic equality checking for results.
func defaultValidateResult(result any, expected map[string]any) error {
	if result == nil && len(expected) > 0 {
		return fmt.Errorf("expected result but got nil")
	}

	if result == nil && len(expected) == 0 {
		return nil // Both nil, considered equal
	}

	// Convert result to map for comparison
	resultMap := make(map[string]any)
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	if err := json.Unmarshal(resultJSON, &resultMap); err != nil {
		return fmt.Errorf("failed to unmarshal result to map: %w", err)
	}

	// Compare each expected field
	for key, expectedValue := range expected {
		actualValue, ok := resultMap[key]
		if !ok {
			return fmt.Errorf("missing expected field %q", key)
		}

		// Convert both to JSON for deep comparison
		expectedJSON, _ := json.Marshal(expectedValue)
		actualJSON, _ := json.Marshal(actualValue)
		if string(expectedJSON) != string(actualJSON) {
			return fmt.Errorf("field %q: expected %s, got %s", key, expectedJSON, actualJSON)
		}
	}

	return nil
}

func (r *ScenarioRunner) validateStream(t *testing.T, method string, stream any, expect Expectation) {
	if stream == nil {
		t.Errorf("stream is nil for method %s", method)
		return
	}

	// Stream validation with custom validators
	if expect.Validator != "" {
		t.Logf("Stream validator %s specified for %s - implement stream validation", expect.Validator, method)
		return
	}

	// No default stream validation - streams are too varied
	t.Logf("Stream validation for %s: specify a validator in YAML or implement custom validation", method)
}

func (r *ScenarioRunner) selectTransport(client *Client, transport string) *Client {
	switch transport {
	case "http", "http-sse", "http-ws":
		return client.HTTP()
	case "grpc":
		return client.GRPC()
	case "jsonrpc", "jsonrpc-sse", "jsonrpc-ws":
		return client.JSONRPC()
	default:
		return client // auto or unknown - use default
	}
}
`

var ScenarioRunnerWithoutResultCode = `// ScenarioRunner executes test scenarios.
type ScenarioRunner struct {
	scenarios  []Scenario
	validators Validators // Global validator configuration
}

// LoadScenarios loads scenarios from a YAML file.
func LoadScenarios(path string) (*ScenarioRunner, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenarios file: %w", err)
	}

	var config ScenarioConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse scenarios YAML: %w", err)
	}

	return &ScenarioRunner{
		scenarios:  config.Scenarios,
		validators: config.Validators,
	}, nil
}

// NewScenarioRunner creates a new scenario runner.
func NewScenarioRunner() *ScenarioRunner {
	return &ScenarioRunner{
		scenarios: make([]Scenario, 0),
	}
}

// AddScenario adds a scenario to the runner.
func (r *ScenarioRunner) AddScenario(scenario Scenario) {
	r.scenarios = append(r.scenarios, scenario)
}

// Run executes all scenarios.
func (r *ScenarioRunner) Run(t *testing.T, client *Client) {
	if r == nil {
		t.Fatal("ScenarioRunner is nil")
	}
	if client == nil {
		t.Fatal("Client is nil")
	}
	for _, scenario := range r.scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			r.runScenario(t, client, scenario)
		})
	}
}

// RunNamed executes a specific scenario by name.
func (r *ScenarioRunner) RunNamed(t *testing.T, client *Client, name string) {
	if r == nil {
		t.Fatal("ScenarioRunner is nil")
	}
	if client == nil {
		t.Fatal("Client is nil")
	}
	if name == "" {
		t.Fatal("scenario name is empty")
	}
	for _, scenario := range r.scenarios {
		if scenario.Name == name {
			r.runScenario(t, client, scenario)
			return
		}
	}
	t.Fatalf("scenario %q not found", name)
}

func (r *ScenarioRunner) runScenario(t *testing.T, client *Client, scenario Scenario) {
	// Apply default transport if specified
	if scenario.Transport != "" {
		client = r.selectTransport(client, scenario.Transport)
	}

	for i, step := range scenario.Steps {
		t.Run(fmt.Sprintf("step_%d_%s", i+1, step.Method), func(t *testing.T) {
			// Apply scenario-level timeout if step doesn't override
			if step.Timeout == "" && scenario.Timeout != "" {
				step.Timeout = scenario.Timeout
			}
			r.runStep(t, client, step)
		})
	}
}

func (r *ScenarioRunner) runStep(t *testing.T, client *Client, step Step) {
	// Apply per-step transport override
	if step.Transport != "" {
		client = r.selectTransport(client, step.Transport)
	}

	// Validate transport availability
	if step.Transport != "" && step.Transport != "auto" {
		if transports, ok := TransportAvailability[step.Method]; ok {
			found := false
			for _, t := range transports {
				if t == step.Transport {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("method %q does not support transport %q, available: %v",
					step.Method, step.Transport, transports)
			}
		}
	}

	// Process payload
	payload := step.Payload
	ctx := context.Background()

	// Apply timeout if specified
	if step.Timeout != "" {
		duration, err := time.ParseDuration(step.Timeout)
		if err != nil {
			t.Fatalf("invalid timeout %q: %v", step.Timeout, err)
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, duration)
		defer cancel()
	}

	// Execute the method
	result, err := r.executeMethod(ctx, client, step.Method, payload)

	// Handle error expectation
	if step.Expect.Error != "" {
		if err == nil {
			t.Errorf("expected error %q but got none", step.Expect.Error)
		} else if !strings.Contains(err.Error(), step.Expect.Error) {
			t.Errorf("expected error containing %q but got %q", step.Expect.Error, err.Error())
		}
		return
	}

	// Handle unexpected error
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	// Validate result if expected
	if step.Expect.Result != nil || step.Expect.Validator != "" {
		r.validateResult(t, step.Method, result, step.Expect)
	}

	// Handle streaming expectations
	if len(step.Expect.Stream) > 0 {
		r.validateStream(t, step.Method, result, step.Expect)
	}
}

func (r *ScenarioRunner) executeMethod(ctx context.Context, client *Client, method string, payload map[string]any) (any, error) {
	switch method {
	case "WithoutResultMethod":
		return nil, client.WithoutResultMethod(ctx)
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

func (r *ScenarioRunner) mapToStruct(data map[string]any, target any) error {
	if data == nil {
		// nil data is okay, just return without setting anything
		return nil
	}
	if target == nil {
		return fmt.Errorf("target is nil")
	}
	// Convert map to JSON then unmarshal to struct
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}

func (r *ScenarioRunner) validateResult(t *testing.T, method string, result any, expect Expectation) {
	if result == nil && expect.Result == nil && expect.Validator == "" {
		// Nothing to validate
		return
	}

	// If custom validator specified in YAML, call it
	if expect.Validator != "" {
		// Call the user-defined validator function
		// The function signature should be: func(t *testing.T, result *ServiceType, expected map[string]any)
		r.callValidator(t, method, result, expect)
		return
	}

	// Fall back to default validation
	if expect.Result != nil {
		if err := defaultValidateResult(result, expect.Result); err != nil {
			t.Errorf("validation failed for %s: %v", method, err)
		}
	}
}

// callValidator calls the user-specified validator function.
// The validator function must be defined in the test package.
func (r *ScenarioRunner) callValidator(t *testing.T, method string, result any, expect Expectation) {
	// For each validator found in YAML, we generate a direct call
	// Users must define these functions in their test files

	validatorName := expect.Validator
	_ = validatorName // avoid unused variable in case no validators are defined

	switch method {
	case "WithoutResultMethod":
		t.Errorf("method %q has no result to validate", method)
	default:
		t.Errorf("unknown method: %s", method)
	}
}

// defaultValidateResult provides basic equality checking for results.
func defaultValidateResult(result any, expected map[string]any) error {
	if result == nil && len(expected) > 0 {
		return fmt.Errorf("expected result but got nil")
	}

	if result == nil && len(expected) == 0 {
		return nil // Both nil, considered equal
	}

	// Convert result to map for comparison
	resultMap := make(map[string]any)
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}
	if err := json.Unmarshal(resultJSON, &resultMap); err != nil {
		return fmt.Errorf("failed to unmarshal result to map: %w", err)
	}

	// Compare each expected field
	for key, expectedValue := range expected {
		actualValue, ok := resultMap[key]
		if !ok {
			return fmt.Errorf("missing expected field %q", key)
		}

		// Convert both to JSON for deep comparison
		expectedJSON, _ := json.Marshal(expectedValue)
		actualJSON, _ := json.Marshal(actualValue)
		if string(expectedJSON) != string(actualJSON) {
			return fmt.Errorf("field %q: expected %s, got %s", key, expectedJSON, actualJSON)
		}
	}

	return nil
}

func (r *ScenarioRunner) validateStream(t *testing.T, method string, stream any, expect Expectation) {
	if stream == nil {
		t.Errorf("stream is nil for method %s", method)
		return
	}

	// Stream validation with custom validators
	if expect.Validator != "" {
		t.Logf("Stream validator %s specified for %s - implement stream validation", expect.Validator, method)
		return
	}

	// No default stream validation - streams are too varied
	t.Logf("Stream validation for %s: specify a validator in YAML or implement custom validation", method)
}

func (r *ScenarioRunner) selectTransport(client *Client, transport string) *Client {
	switch transport {
	case "http", "http-sse", "http-ws":
		return client.HTTP()
	case "grpc":
		return client.GRPC()
	case "jsonrpc", "jsonrpc-sse", "jsonrpc-ws":
		return client.JSONRPC()
	default:
		return client // auto or unknown - use default
	}
}
`

var SuiteTestWithResultCode = `// RunWithResultServiceHarness exercises the generated harness against your
// service implementation.// Call this helper from your test, passing your service implementation.
func RunWithResultServiceHarness(t *testing.T, svc withresultservice.Service) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h := withResultServicetest.NewHarness(t, svc)
	defer h.Close()

	td := withResultServicetest.NewTestData()
	t.Run("WithResultMethod", func(t *testing.T) {
		result, err := h.Client.WithResultMethod(ctx)
		if err != nil {
			t.Errorf("WithResultMethod failed: %v", err)
		}
		if result == nil {
			t.Error("WithResultMethod returned nil result")
		}
	})
}
`

var SuiteTestWithoutResultCode = `// RunWithoutResultServiceHarness exercises the generated harness against your
// service implementation.// Call this helper from your test, passing your service implementation.
func RunWithoutResultServiceHarness(t *testing.T, svc withoutresultservice.Service) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	h := withoutResultServicetest.NewHarness(t, svc)
	defer h.Close()

	td := withoutResultServicetest.NewTestData()
	t.Run("WithoutResultMethod", func(t *testing.T) {
		result, err := h.Client.WithoutResultMethod(ctx)
		if err != nil {
			t.Errorf("WithoutResultMethod failed: %v", err)
		}
	})
}
`
