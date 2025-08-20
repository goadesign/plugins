# Example test scenarios for {{ .Name }} service
# This file demonstrates the YAML scenario testing capability.
# Customize these scenarios to match your testing needs.

# Optional: Configure global validator settings
# validators:
#   package: myvalidators  # Package containing validator functions
#   path: myapp/testing/validators  # Import path for the package

scenarios:
{{- if .Methods }}
  # Basic CRUD lifecycle test
  - name: "basic_lifecycle"
    description: "Tests basic service functionality"
    {{- if .HasHTTP }}
    transport: http  # Optional: specify default transport (http, grpc, jsonrpc)
    {{- else if .HasGRPC }}
    transport: grpc
    {{- else if .HasJSONRPC }}
    transport: jsonrpc
    {{- end }}
    steps:
    {{- range $i, $method := .Methods }}
    {{- if eq $i 0 }}
      # Example step for {{ $method.Name }} method
      - method: {{ $method.Name }}
        {{- if $method.Payload }}
        payload:
          # Add your test payload here
          # Example: name: "test"
        {{- end }}
        {{- if $method.Result }}
        expect:
          result:
            # Add expected result fields here
            # Example: status: "success"
          # Optional: specify custom validator function
          # validator: ValidateCustomResult  # Function name in current package
          # validator_pkg: myvalidators     # Or specify different package
        {{- end }}
    {{- end }}
    {{- end }}

  # Transport-specific testing
  {{- if .HasHTTP }}
  - name: "http_specific"
    description: "Tests HTTP-specific behavior"
    transport: http
    steps:
    {{- range $method := .Methods }}
    {{- range $transport := $method.Transports }}
    {{- if eq $transport "http" }}
      - method: {{ $method.Name }}
        {{- if $method.Payload }}
        payload: {}
        {{- end }}
        {{- if $method.Result }}
        expect:
          result: {}
        {{- end }}
        {{- break }}
    {{- end }}
    {{- end }}
    {{- end }}
  {{- end }}

  {{- if .HasGRPC }}
  - name: "grpc_specific"
    description: "Tests gRPC-specific behavior"
    transport: grpc
    steps:
    {{- range $method := .Methods }}
    {{- range $transport := $method.Transports }}
    {{- if eq $transport "grpc" }}
      - method: {{ $method.Name }}
        {{- if $method.Payload }}
        payload: {}
        {{- end }}
        {{- if $method.Result }}
        expect:
          result: {}
        {{- end }}
        {{- break }}
    {{- end }}
    {{- end }}
    {{- end }}
  {{- end }}

  # Error handling test
  - name: "error_handling"
    description: "Tests error conditions"
    steps:
    {{- range $method := .Methods }}
    {{- if $method.Errors }}
    {{- if eq (len $method.Errors) 0 }}
      # {{ $method.Name }} has error handling
      - method: {{ $method.Name }}
        {{- if $method.Payload }}
        payload:
          # Add invalid payload to trigger error
        {{- end }}
        expect:
          error: "expected_error"  # Replace with actual error name
        {{- break }}
    {{- end }}
    {{- end }}
    {{- end }}

  {{- range $method := .Methods }}
  {{- if or (eq $method.StreamKind 1) (eq $method.StreamKind 2) (eq $method.StreamKind 3) }}
  # Streaming test example
  - name: "streaming_{{ $method.Name }}"
    description: "Tests streaming for {{ $method.Name }}"
    steps:
      - method: {{ $method.Name }}
        {{- if eq $method.StreamKind 1 }}
        # Client streaming
        send:
          - # First message
          - # Second message
        expect:
          result: {}
        {{- else if eq $method.StreamKind 2 }}
        # Server streaming
        {{- if $method.Payload }}
        payload: {}
        {{- end }}
        receive:
          - # Expected first message
          - # Expected second message
        {{- else if eq $method.StreamKind 3 }}
        # Bidirectional streaming
        send:
          - # First client message
          - # Second client message
        receive:
          - # Expected first server message
          - # Expected second server message
        {{- end }}
        {{- break }}
  {{- end }}
  {{- end }}
{{- else }}
  # No methods found in service
  - name: "placeholder"
    description: "Add your test scenarios here"
    steps:
      - method: "YourMethod"
        payload: {}
        expect:
          result: {}
{{- end }}

# Note: The scenario runner performs basic smoke testing.
# For detailed assertions and complex test logic, write custom Go tests.

# Transport values (based on your service configuration):
# {{- range $transport := .ValidTransports }}
# - {{ $transport }}
# {{- end }}

# Available methods and their transports:
# {{- range $method := .Methods }}
# - {{ $method.Name }}: {{ range $i, $t := $method.Transports }}{{- if $i }}, {{ end }}{{ $t }}{{- end }}
# {{- end }}