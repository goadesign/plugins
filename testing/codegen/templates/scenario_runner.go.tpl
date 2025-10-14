{{ printf "ScenarioRunner executes test scenarios." | comment }}
type ScenarioRunner struct {
	scenarios []Scenario
	validators Validators // Global validator configuration
}

{{ printf "LoadScenarios loads scenarios from a YAML file." | comment }}
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
		scenarios: config.Scenarios,
		validators: config.Validators,
	}, nil
}

{{ printf "NewScenarioRunner creates a new scenario runner." | comment }}
func NewScenarioRunner() *ScenarioRunner {
	return &ScenarioRunner{
		scenarios: make([]Scenario, 0),
	}
}

{{ printf "AddScenario adds a scenario to the runner." | comment }}
func (r *ScenarioRunner) AddScenario(scenario Scenario) {
	r.scenarios = append(r.scenarios, scenario)
}

{{ printf "Run executes all scenarios." | comment }}
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

{{ printf "RunNamed executes a specific scenario by name." | comment }}
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
	{{- range .Methods }}
	case "{{ .Name }}":
		{{- if .PayloadRef }}
		// Convert payload map to typed payload
		p := &{{ $.PkgName }}.{{ .Payload }}{}
		if err := r.mapToStruct(payload, p); err != nil {
			return nil, fmt.Errorf("invalid payload for {{ .Name }}: %w", err)
		}
		{{- end }}
		return client.{{ .VarName }}(ctx{{ if .PayloadRef }}, p{{ end }})
	{{- end }}
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

{{ printf "callValidator calls the user-specified validator function." | comment }}
{{ printf "The validator function must be defined in the test package." | comment }}
func (r *ScenarioRunner) callValidator(t *testing.T, method string, result any, expect Expectation) {
	// For each validator found in YAML, we generate a direct call
	// Users must define these functions in their test files
	
	validatorName := expect.Validator
	_ = validatorName // avoid unused variable in case no validators are defined
	
	switch method {
	{{- range .Methods }}
	case "{{ .Name }}":
		{{- if .ResultRef }}
		typedResult := result.(*{{ $.PkgName }}.{{ .Result }})
		{{- $validators := index $.Validators .Name }}
		{{- if $validators }}
		
		// Call the validator function by name
		switch validatorName {
		{{- range $validators }}
		case "{{ . }}":
			// Direct function call - user must define this function
			{{- if $.ValidatorPkg }}
			{{ $.ValidatorPkg }}.{{ . }}(t, typedResult, expect.Result)
			return
			{{- else }}
			{{ . }}(t, typedResult, expect.Result)
			return
			{{- end }}
		{{- end }}
		default:
			t.Errorf("unknown validator %q for method %q", validatorName, method)
		}
		{{- else }}
		_ = typedResult // no validators defined in YAML
		t.Errorf("validator %q specified but not generated - add it to scenarios.yaml first", validatorName)
		{{- end }}
		{{- else }}
		t.Errorf("method %q has no result to validate", method)
		{{- end }}
	{{- end }}
	default:
		t.Errorf("unknown method: %s", method)
	}
}

{{ printf "defaultValidateResult provides basic equality checking for results." | comment }}
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
	{{- if .HasHTTP }}
	case "http", "http-sse", "http-ws":
		return client.HTTP()
	{{- end }}
	{{- if .HasGRPC }}
	case "grpc":
		return client.GRPC()
	{{- end }}
	{{- if .HasJSONRPC }}
	case "jsonrpc", "jsonrpc-sse", "jsonrpc-ws":
		return client.JSONRPC()
	{{- end }}
	default:
		return client // auto or unknown - use default
	}
}