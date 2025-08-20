# Goa Testing Plugin

> **Transform your Goa services into battle-tested, production-ready APIs with comprehensive transport-aware testing.**

The Goa Testing Plugin generates a complete testing framework that validates
your services through real transport layers (HTTP, gRPC, JSON-RPC),
automatically verifying every response against your DSL contracts. Write tests
once, run them across all transports, and catch integration issues before they
reach production.

## Table of Contents

- [Why This Plugin?](#why-this-plugin)
- [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
- [Generated Test Suite](#generated-test-suite)
- [YAML Scenario Testing](#yaml-scenario-testing)
- [Custom Validators](#custom-validators)
- [Generated Test Helpers](#generated-test-helpers)
- [Streaming Support](#streaming-support)
- [Transport-Specific Testing](#transport-specific-testing)
- [Advanced Features](#advanced-features)
- [Best Practices](#best-practices)
- [File Organization](#file-organization)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)

## Why This Plugin?

Traditional unit tests mock transport layers and miss critical integration bugs.
This plugin revolutionizes service testing by:

### **Real Transport Validation**
Your tests run through actual HTTP servers, gRPC connections, and WebSocket
streams - exactly as they will in production. No more "works in unit tests,
fails in production" surprises.

### **Automatic Contract Enforcement**
Every response is automatically validated against your DSL. Required fields,
types, formats, and validation rules—if it's defined in your DSL, it's tested.
No manual assertions are needed for structural validation. While Goa ensures
incoming requests are valid, this plugin extends that rigor to response contents
as well.

### **Generated Test Suite**

The plugin generates a comprehensive test suite tailored to your service,
covering all endpoints and transports defined in your DSL. This suite is ready
to use out of the box and exercises your service through real HTTP, gRPC, and
other supported transports, ensuring end-to-end contract compliance.

### **Declarative Scenario Testing**
Define complex test scenarios in YAML. Perfect for regression tests, acceptance
criteria, and API documentation that doubles as executable tests.

## Quick Start

### 1. Installation

Add the plugin to your design:

```go
package design

import (
    . "goa.design/goa/v3/dsl"
    _ "goa.design/plugins/v3/testing"  // Add this line
)
```

### 2. Generate Testing Framework

```bash
# Generate test harness, client, and helpers (required)
$ goa gen your/design/package

# Generate editable test suite scaffold (optional, recommended)
$ goa example your/design/package
```

The `goa example` command creates a test suite file that won't be overwritten on
subsequent runs - it's yours to customize.

### 3. Write Your First Test

Using the generated test suite:

```go
// calculator_test.go
func TestCalculator(t *testing.T) {
    svc := NewCalculatorService()  // Your implementation
    RunCalculatorHarness(t, svc)   // Generated test suite
}
```

Or write custom tests:

```go
func TestCustomScenario(t *testing.T) {
    svc := NewCalculatorService()
    h := calculatortest.NewHarness(t, svc)
    defer h.Close()
    
    // Every call validates against your DSL automatically
    result, err := h.Client.Add(ctx, &AddPayload{A: 2, B: 3})
    require.NoError(t, err)
    assert.Equal(t, 5, result.Sum)
}
```

## Core Concepts

### The Test Harness

The harness is your testing command center, managing test servers and providing
transport-aware clients:

```go
h := myservicetest.NewHarness(t, service)
defer h.Close()
```

Behind the scenes, the harness:
- Starts an `httptest.Server` for HTTP endpoints
- Creates an in-memory gRPC server using `bufconn`
- Configures WebSocket upgraders for streaming
- Sets up JSON-RPC handlers
- Manages all cleanup automatically

### Transport-Aware Client

The test client adapts to any transport, letting you test the same logic across different protocols:

```go
// Auto-select transport based on availability
result, err := h.Client.GetUser(ctx, payload)

// Explicitly test gRPC behavior
result, err := h.Client.GRPC().GetUser(ctx, payload)

// Test HTTP-specific features
result, err := h.Client.HTTP().GetUser(ctx, payload)

// Test streaming variants
stream, err := h.Client.HTTP().AsStream().Subscribe(ctx)
```

### Automatic DSL Validation

Every response is validated automatically. No manual assertions needed for DSL contracts:

```go
// Your DSL defines:
Field("email", String, func() {
    Format(FormatEmail)
    Pattern(`^[a-z]+@[a-z]+\.com$`)
    Required()
})

// If service returns: {"email": "INVALID@EXAMPLE.COM"}
// Test fails with: "email: does not match pattern '^[a-z]+@[a-z]+\.com$'"
```

## Generated Test Suite

The `goa example` command generates a complete test suite scaffold that
exercises every service method. This is editable code - customize it freely!

### Generated Structure

```go
// calculator_suite_test.go (generated once, never overwritten)
func RunCalculatorHarness(t *testing.T, svc calculator.Service) {
    t.Helper()
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    h := calculatortest.NewHarness(t, svc)
    defer h.Close()
    
    td := calculatortest.NewTestData()
    
    t.Run("Add", func(t *testing.T) {
        result, err := h.Client.Add(ctx, td.ValidAddPayload())
        if err != nil {
            t.Fatalf("Add failed: %v", err)
        }
        if result == nil {
            t.Error("Add returned nil result")
        }
    })
    
    t.Run("Subscribe_Stream", func(t *testing.T) {
        stream, err := h.Client.Subscribe(ctx)
        // ... streaming test logic
    })
}
```

### Features Included

- **Helper Function Pattern**: `RunXXXHarness` accepts your service implementation
- **Subtests**: Each method gets its own `t.Run()` for isolated testing
- **Valid Test Data**: Uses generated data builders for valid payloads
- **Timeout Management**: Configurable context timeouts for all operations
- **Proper Cleanup**: Deferred cleanup ensures resources are freed
- **Streaming Support**: Handles all streaming patterns correctly

## YAML Scenario Testing

Define test scenarios declaratively in YAML for readable, maintainable tests that non-developers can understand and review.

### Complete YAML Schema

```yaml
# scenarios.yaml
validators:                         # Optional: Custom validator configuration
  package: "testvalidators"        # Package name for validators
  path: "myapp/testvalidators"     # Import path for validator package

scenarios:                          # Required: List of test scenarios
  - name: "user_lifecycle"         # Required: Unique scenario identifier
    description: "Test CRUD ops"   # Optional: Human-readable description
    transport: "grpc"              # Optional: Default transport for all steps
    timeout: "30s"                 # Optional: Default timeout for all steps
    steps:                         # Required: List of test steps
      - method: "CreateUser"       # Required: Service method name
        transport: "http"          # Optional: Override scenario transport
        timeout: "5s"              # Optional: Override scenario timeout
        payload:                   # Optional: Input data (matches DSL structure)
          name: "Alice"
          email: "alice@example.com"
          age: 30
        stream: false              # Optional: Use streaming variant (default: false)
        send:                      # Optional: For client/bidi streaming
          - message: "Hello"
          - message: "World"
        receive:                   # Optional: Expected stream messages
          - message: "Got: Hello"
          - message: "Got: World"
        expect:                    # Optional: Expected outcome
          result:                  # Optional: Expected response fields
            id: "user-123"
            status: "created"
          error: ""                # Optional: Expected error substring
          stream:                  # Optional: Expected stream messages
            - data: "update1"
            - data: "update2"
          validator: "ValidateUserCreation"  # Optional: Custom validator function
```

### Field Reference

#### Scenario Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Unique identifier for the scenario |
| `description` | string | ❌ | Human-readable description for documentation |
| `transport` | string | ❌ | Default transport: `auto`, `http`, `grpc`, `jsonrpc`, `http-ws`, `http-sse`, `jsonrpc-ws`, `jsonrpc-sse` |
| `timeout` | duration | ❌ | Default timeout for all steps (e.g., `10s`, `500ms`, `1m`) |
| `steps` | []Step | ✅ | Ordered list of test steps to execute |

#### Step Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `method` | string | ✅ | Service method name as defined in DSL |
| `transport` | string | ❌ | Override scenario's default transport |
| `timeout` | duration | ❌ | Override scenario's default timeout |
| `payload` | object | ❌ | Input data matching your DSL payload structure |
| `stream` | bool | ❌ | Use streaming variant for mixed endpoints |
| `send` | []object | ❌ | Messages to send (client/bidi streaming) |
| `receive` | []object | ❌ | Expected received messages (server/bidi streaming) |
| `expect` | Expectation | ❌ | Expected outcome (see below) |

#### Expectation Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `result` | object | ❌ | Expected response fields (partial match) |
| `error` | string | ❌ | Expected error message substring |
| `stream` | []object | ❌ | Expected stream messages in order |
| `validator` | string | ❌ | Custom validator function name |

### Transport Options Explained

- **`auto`**: Use first available transport (default)
- **`http`**: Plain HTTP (REST endpoints)
- **`grpc`**: gRPC unary or streaming
- **`jsonrpc`**: JSON-RPC over HTTP
- **`http-ws`**: HTTP WebSocket (for streaming)
- **`http-sse`**: HTTP Server-Sent Events
- **`jsonrpc-ws`**: JSON-RPC over WebSocket
- **`jsonrpc-sse`**: JSON-RPC over SSE

### Running YAML Scenarios

```go
func TestScenarios(t *testing.T) {
    service := NewMyService()
    h := myservicetest.NewHarness(t, service)
    defer h.Close()
    
    // Load and run all scenarios
    runner, err := myservicetest.LoadScenarios("scenarios.yaml")
    require.NoError(t, err)
    runner.Run(t, h.Client)
    
    // Or run specific scenario
    runner.RunNamed(t, h.Client, "user_lifecycle")
}
```

### Real-World Examples

#### Testing User Registration Flow

```yaml
scenarios:
  - name: "user_registration"
    description: "Complete user registration with validation"
    timeout: "10s"
    steps:
      - method: "CheckEmailAvailable"
        payload:
          email: "newuser@example.com"
        expect:
          result:
            available: true
            
      - method: "CreateUser"
        payload:
          email: "newuser@example.com"
          password: "SecurePass123!"
          name: "New User"
        expect:
          result:
            status: "pending_verification"
          validator: "ValidateUserCreation"
          
      - method: "SendVerificationEmail"
        payload:
          user_id: "12345"  # Use known test ID
        expect:
          result:
            sent: true
```

#### Testing Error Conditions

```yaml
scenarios:
  - name: "error_handling"
    description: "Verify proper error responses"
    steps:
      - method: "Divide"
        payload:
          numerator: 10
          divisor: 0
        expect:
          error: "division by zero"
          
      - method: "GetUser"
        payload:
          id: "nonexistent"
        expect:
          error: "user not found"
```

#### Testing Streaming

```yaml
scenarios:
  - name: "live_updates"
    description: "Test real-time update stream"
    transport: "http-ws"
    steps:
      - method: "SubscribeToUpdates"
        payload:
          topics: ["news", "alerts"]
        expect:
          stream:
            - type: "news"
              data: "Breaking news"
            - type: "alert"
              data: "System update"
```

## Custom Validators

While Goa validates structural contracts automatically, custom validators verify business logic, calculations, and complex rules.

### Defining Validators

**Important**: Validators must be in a separate package to avoid import cycles.

```yaml
# scenarios.yaml
validators:
  package: "validators"
  path: "myapp/validators"

scenarios:
  - name: "pricing_test"
    steps:
      - method: "CalculatePrice"
        payload:
          items: [...]
        expect:
          validator: "ValidatePricing"
```

```go
// validators/pricing.go
package validators

import (
    "testing"
    "myapp/gen/pricing"
)

// Validator signature: func(t *testing.T, result *TypedResult, expected map[string]any)
func ValidatePricing(t *testing.T, result *pricing.Quote, expected map[string]any) {
    // Validate business rules
    if result.Discount > result.Subtotal {
        t.Errorf("discount ($%.2f) exceeds subtotal ($%.2f)", 
            result.Discount, result.Subtotal)
    }
    
    // Validate calculations
    expectedTotal := result.Subtotal - result.Discount + result.Tax
    if math.Abs(result.Total - expectedTotal) > 0.01 {
        t.Errorf("total calculation error: expected %.2f, got %.2f", 
            expectedTotal, result.Total)
    }
    
    // Use expected values from YAML if needed
    if maxDiscount, ok := expected["max_discount"].(float64); ok {
        if result.Discount > maxDiscount {
            t.Errorf("discount exceeds maximum: %.2f > %.2f", 
                result.Discount, maxDiscount)
        }
    }
}
```

### Common Validator Patterns

#### Timing Validation

```go
func ValidateResponseTime(t *testing.T, result *Response, expected map[string]any) {
    maxMs := 100
    if result.ResponseTimeMs > maxMs {
        t.Errorf("response too slow: %dms (max: %dms)", 
            result.ResponseTimeMs, maxMs)
    }
}
```

#### Data Consistency

```go
func ValidateDataIntegrity(t *testing.T, result *Report, expected map[string]any) {
    sum := 0
    for _, item := range result.Items {
        sum += item.Count
    }
    if sum != result.TotalCount {
        t.Errorf("data inconsistency: sum of items (%d) != total (%d)", 
            sum, result.TotalCount)
    }
}
```

#### Complex State Validation

```go
func ValidateStateTransition(t *testing.T, result *Order, expected map[string]any) {
    validTransitions := map[string][]string{
        "pending": {"confirmed", "cancelled"},
        "confirmed": {"shipped", "cancelled"},
        "shipped": {"delivered", "returned"},
    }
    
    oldState := expected["old_state"].(string)
    valid := validTransitions[oldState]
    
    found := false
    for _, v := range valid {
        if v == result.Status {
            found = true
            break
        }
    }
    
    if !found {
        t.Errorf("invalid state transition: %s -> %s", oldState, result.Status)
    }
}
```

## Generated Test Helpers

The plugin generates comprehensive test utilities to eliminate boilerplate and ensure correctness.

### Test Data Generators

Generate valid test data that satisfies all DSL constraints:

```go
td := myservicetest.NewTestData()

// Returns a fully valid payload with all required fields
payload := td.ValidCreateUserPayload()
// Generates: {
//   "name": "generated-name-x7k2m",
//   "email": "user-9f3kd@example.com",
//   "age": 25,
//   "roles": ["user"]
// }

// Customize specific fields
payload.Email = "custom@example.com"
payload.Roles = []string{"admin", "user"}
```

Generated data respects all DSL constraints:
- Required fields are always populated
- Formats (email, UUID, etc.) are valid
- Patterns are matched
- Min/max constraints are satisfied
- Enums use valid values

### Error Assertions

Type-safe error checking for DSL-defined errors:

```go
// Your DSL defines:
Error("insufficient_funds", ErrorResult, "Not enough balance")

// Generated assertion:
_, err := h.Client.Transfer(ctx, payload)
myservicetest.AssertInsufficientFunds(t, err)
// Validates error type, message, and any error fields
```

### Transport Error Handling

```go
// Check HTTP status codes
if httpErr, ok := err.(*goa.ServiceError); ok {
    assert.Equal(t, 429, httpErr.StatusCode)  // Rate limited
}

// Check gRPC status codes
if st, ok := status.FromError(err); ok {
    assert.Equal(t, codes.NotFound, st.Code())
}
```

## Streaming Support

Complete support for all streaming patterns with proper lifecycle management.

### Server Streaming

```go
stream, err := h.Client.GetUpdates(ctx, &GetUpdatesPayload{
    Since: time.Now().Add(-1 * time.Hour),
})
require.NoError(t, err)

updates := []Update{}
for {
    update, err := stream.Recv()
    if err == io.EOF {
        break  // Stream ended normally
    }
    require.NoError(t, err)
    updates = append(updates, update)
}

assert.Len(t, updates, 10)  // Validate received data
```

### Client Streaming

```go
stream, err := h.Client.UploadData(ctx)
require.NoError(t, err)

// Send multiple messages
testData := []Data{
    {ID: 1, Value: "first"},
    {ID: 2, Value: "second"},
    {ID: 3, Value: "third"},
}

for _, data := range testData {
    err := stream.Send(&data)
    require.NoError(t, err)
}

// Close and get result
result, err := stream.CloseAndRecv()
require.NoError(t, err)
assert.Equal(t, 3, result.ProcessedCount)
```

### Bidirectional Streaming

```go
stream, err := h.Client.Chat(ctx)
require.NoError(t, err)

// Send messages in goroutine
go func() {
    messages := []string{"Hello", "How are you?", "Goodbye"}
    for _, msg := range messages {
        stream.Send(&ChatMessage{Text: msg})
        time.Sleep(100 * time.Millisecond)
    }
    stream.CloseSend()
}()

// Receive responses
responses := []string{}
for {
    resp, err := stream.Recv()
    if err == io.EOF {
        break
    }
    require.NoError(t, err)
    responses = append(responses, resp.Text)
}

assert.Contains(t, responses, "Hello to you too!")
```

### YAML Streaming Tests

```yaml
scenarios:
  - name: "streaming_chat"
    transport: "grpc"
    steps:
      - method: "Chat"
        send:
          - text: "Hello"
          - text: "How are you?"
        receive:
          - text: "Hi there!"
          - text: "I'm doing well, thanks!"
        expect:
          validator: "ValidateChatSession"
```

## Transport-Specific Testing

Test transport-specific features and behaviors while maintaining the same API.

### HTTP-Specific Features

```go
// Test custom headers
req := h.HTTPRequest("GET", "/users/123", nil)
req.Header.Set("X-Custom-Header", "value")
resp := h.HTTPDo(req)
assert.Equal(t, 200, resp.StatusCode)

// Test SSE (Server-Sent Events)
stream, err := h.Client.HTTP().AsStream().Subscribe(ctx, payload)
require.NoError(t, err)

event := <-stream.Events()
assert.Equal(t, "update", event.Type)

// Test WebSocket
wsURL := h.HTTPWSURL("/chat")
conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
require.NoError(t, err)
defer conn.Close()
```

### gRPC-Specific Features

```go
// Test with metadata
import "google.golang.org/grpc/metadata"

ctx = metadata.AppendToOutgoingContext(ctx, 
    "api-key", "secret",
    "trace-id", "12345",
)
result, err := h.Client.GRPC().GetUser(ctx, payload)

// Verify gRPC status codes
st, ok := status.FromError(err)
require.True(t, ok)
assert.Equal(t, codes.PermissionDenied, st.Code())
assert.Contains(t, st.Message(), "insufficient permissions")

// Test gRPC interceptors
conn := h.GRPCConn()
// Add your interceptors for testing
```

### JSON-RPC Features

```go
// Test batch requests
batch := []jsonrpc.Request{
    {Method: "GetUser", Params: map[string]any{"id": "1"}},
    {Method: "GetUser", Params: map[string]any{"id": "2"}},
}
results, err := h.Client.JSONRPC().Batch(ctx, batch)

// Test notifications (no response expected)
err = h.Client.JSONRPC().Notify(ctx, "LogEvent", params)
assert.NoError(t, err)  // No response to validate

// Test over WebSocket
stream, err := h.Client.JSONRPC().AsStream().Connect(ctx)
```

## Advanced Features

### Parallel Testing

The harness is safe for parallel execution:

```go
func TestParallel(t *testing.T) {
    t.Parallel()  // Safe!
    
    service := NewService()
    h := myservicetest.NewHarness(t, service)
    defer h.Close()
    
    // Each harness gets isolated servers
    // No port conflicts or shared state
}
```

### Timeout Configuration

Control timeouts at multiple levels:

```go
// Harness-level timeout
h := myservicetest.NewHarness(t, service, 
    myservicetest.WithTimeout(30 * time.Second))

// Context timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
result, err := h.Client.SlowMethod(ctx, payload)

// YAML timeout
scenarios:
  - name: "performance"
    timeout: "1m"  # Scenario default
    steps:
      - method: "QuickOp"
        timeout: "100ms"  # Tight timeout for fast operations
```

### Custom Transport Selection

```go
// Force specific transport for testing
client := h.Client.GRPC()  // Always use gRPC
client := h.Client.HTTP()   // Always use HTTP

// Check transport availability
if h.Client.HasGRPC() {
    // Run gRPC-specific tests
}

// Dynamic transport selection
transport := os.Getenv("TEST_TRANSPORT")
switch transport {
case "grpc":
    client = h.Client.GRPC()
case "http":
    client = h.Client.HTTP()
default:
    client = h.Client  // Auto-select
}
```

### Mock Service Integration

```go
type MockService struct {
    // Embed the interface for easy partial mocking
    calculator.Service
    
    AddFunc func(context.Context, *AddPayload) (*AddResult, error)
}

func (m *MockService) Add(ctx context.Context, p *AddPayload) (*AddResult, error) {
    if m.AddFunc != nil {
        return m.AddFunc(ctx, p)
    }
    // Default behavior
    return &AddResult{Sum: p.A + p.B}, nil
}

func TestWithMock(t *testing.T) {
    mock := &MockService{
        AddFunc: func(ctx context.Context, p *AddPayload) (*AddResult, error) {
            // Custom mock behavior
            return nil, errors.New("simulated error")
        },
    }
    
    h := calculatortest.NewHarness(t, mock)
    defer h.Close()
    
    _, err := h.Client.Add(ctx, payload)
    assert.Error(t, err)
}
```

## Best Practices

### 1. Use the Generated Test Suite

Start with the scaffold - it's a complete working example:

```bash
$ goa example your/design/package  # Generate once
$ go test                          # Run immediately
```

The generated suite provides proper setup, cleanup, and test structure.

### 2. Test Through Transports

Always test through the harness client:

```go
// ❌ Don't: Direct service calls miss transport validation
result, err := service.GetUser(ctx, payload)

// ✅ Do: Transport testing catches serialization issues
result, err := h.Client.GetUser(ctx, payload)
```

### 3. Leverage Test Data Generators

Start with valid generated data:

```go
// ❌ Don't: Manual construction might miss required fields
payload := &UserPayload{Name: "John"}  // Oops, forgot required email!

// ✅ Do: Generated data is always valid
payload := td.ValidUserPayload()
payload.Name = "John"  // Customize what you need
```

### 4. YAML for Regression Tests

Capture bugs as scenarios:

```yaml
scenarios:
  - name: "regression_issue_451"
    description: "Ensure timezone handling is correct"
    steps:
      - method: "ScheduleEvent"
        payload:
          time: "2024-03-10T02:30:00Z"  # DST transition
          timezone: "America/New_York"
        expect:
          result:
            local_time: "2024-03-09T21:30:00-05:00"
```

### 5. Separate Validators Package

Always put validators in a separate package:

```
myapp/
├── design/
├── gen/
│   └── myservice/
│       └── myservicetest/     # Generated test package
├── validators/                 # Your validators here
│   └── validators.go          # Separate package avoids cycles
└── myservice_test.go          # Your tests
```

### 6. Test Error Paths

Don't just test the happy path:

```go
func TestErrorHandling(t *testing.T) {
    testCases := []struct {
        name    string
        payload *Payload
        errMsg  string
    }{
        {"missing required", &Payload{}, "name is required"},
        {"invalid email", &Payload{Name: "x", Email: "bad"}, "invalid email"},
        {"negative age", &Payload{Name: "x", Email: "x@y.z", Age: -1}, "age must be positive"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := h.Client.Create(ctx, tc.payload)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tc.errMsg)
        })
    }
}
```

## File Organization

Understanding the generated file structure helps you navigate and extend tests effectively:

```
your-project/
├── design/
│   └── design.go                # Your API design with testing plugin
├── gen/                         # Generated by 'goa gen' (don't edit)
│   └── myservice/
│       ├── service.go           # Service interfaces
│       ├── client.go            # Client implementation
│       └── myservicetest/       # Testing package
│           ├── client.go        # Test client with transport selection
│           ├── harness.go       # Test harness managing servers
│           ├── scenarios.go     # YAML scenario runner
│           ├── testdata.go      # Valid data generators
│           └── errors.go        # Error assertions
├── myservice_suite_test.go      # Generated by 'goa example' (editable)
├── myservice_test.go            # Your custom tests
├── scenarios.yaml               # Your test scenarios
└── validators/                  # Your custom validators
    └── validators.go            # Must be separate package
```

Key points:
- `gen/` is regenerated - don't edit
- `*_suite_test.go` is generated once - customize freely
- Keep validators in a separate package
- Organize scenarios by feature or flow

## Troubleshooting

### Import Cycles

**Problem**: `import cycle not allowed`

**Solution**: Move validators to a separate package:

```go
// ❌ Don't: Validators in test file cause cycles
// myservice_test.go
package myservice_test

import "gen/myservice/myservicetest"  // Test imports this
func ValidateCustom(...) { }           // This tries to use it

// ✅ Do: Separate validators package
// validators/custom.go
package validators

import "gen/myservice"                 // One-way import
func ValidateCustom(...) { }
```

### Missing Methods in Test Suite

**Problem**: Generated test suite doesn't include all methods

**Solution**: Check your DSL - methods need transport bindings:

```go
// ❌ Method without transport won't be tested
Method("Internal", func() {
    Payload(String)
    Result(String)
})

// ✅ Add transport binding
Method("Internal", func() {
    Payload(String)
    Result(String)
    HTTP(func() { POST("/internal") })
    GRPC(func() {})
})
```

### Validator Not Found

**Problem**: `validator "ValidateFoo" not found`

**Solution**: Ensure validator is defined in YAML first:

```yaml
validators:
  package: "validators"
  path: "myapp/validators"

scenarios:
  - steps:
    - expect:
        validator: "ValidateFoo"  # Must match function name exactly
```

### Transport Not Available

**Problem**: `transport "http" not available for method "Foo"`

**Solution**: Check the generated `TransportAvailability` map or use `auto`:

```go
// Check what's available
fmt.Println(myservicetest.TransportAvailability["Foo"])
// Output: ["grpc", "http-ws"]

// Use auto to select first available
transport: "auto"
```

## Examples

### Complete Service Test

```go
func TestUserService(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    service := NewUserService(db)
    h := usertest.NewHarness(t, service)
    defer h.Close()
    
    // Test user lifecycle
    t.Run("CreateUser", func(t *testing.T) {
        td := usertest.NewTestData()
        payload := td.ValidCreateUserPayload()
        
        user, err := h.Client.CreateUser(ctx, payload)
        require.NoError(t, err)
        assert.NotEmpty(t, user.ID)
        assert.Equal(t, payload.Email, user.Email)
    })
    
    t.Run("GetUser", func(t *testing.T) {
        user, err := h.Client.GetUser(ctx, &GetUserPayload{ID: "123"})
        require.NoError(t, err)
        assert.Equal(t, "123", user.ID)
    })
    
    t.Run("UpdateUser", func(t *testing.T) {
        update := &UpdateUserPayload{
            ID: "123",
            Name: "Updated Name",
        }
        user, err := h.Client.UpdateUser(ctx, update)
        require.NoError(t, err)
        assert.Equal(t, "Updated Name", user.Name)
    })
    
    t.Run("DeleteUser", func(t *testing.T) {
        err := h.Client.DeleteUser(ctx, &DeleteUserPayload{ID: "123"})
        require.NoError(t, err)
        
        // Verify deletion
        _, err = h.Client.GetUser(ctx, &GetUserPayload{ID: "123"})
        usertest.AssertUserNotFound(t, err)
    })
}
```

### Performance Testing

```go
func BenchmarkAPI(b *testing.B) {
    service := NewService()
    h := servicetest.NewHarness(b, service)
    defer h.Close()
    
    td := servicetest.NewTestData()
    payload := td.ValidPayload()
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            _, err := h.Client.Process(context.Background(), payload)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

### Integration Test with External Services

```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup external dependencies
    redis := startRedis(t)
    defer redis.Close()
    
    postgres := startPostgres(t)
    defer postgres.Close()
    
    // Create service with real dependencies
    service := NewService(
        WithRedis(redis.URL),
        WithPostgres(postgres.URL),
    )
    
    // Run through test harness
    h := servicetest.NewHarness(t, service)
    defer h.Close()
    
    // Load and run integration scenarios
    runner, err := servicetest.LoadScenarios("integration_scenarios.yaml")
    require.NoError(t, err)
    runner.Run(t, h.Client)
}
```

## Summary

The Goa Testing Plugin transforms testing from a chore into a confidence-building powerhouse:

✨ **Complete test infrastructure** generated automatically  
🎯 **Transport-aware testing** catches real integration issues  
✅ **Automatic DSL validation** ensures contract compliance  
📝 **Declarative YAML scenarios** for maintainable test suites  
🔧 **Rich test helpers** eliminate boilerplate  
🚀 **Zero-friction setup** with generated test suites  

Stop writing test infrastructure. Start testing what matters. Your services, validated through real transports, with complete confidence.

## Getting Help

- 📖 [Goa Documentation](https://goa.design)
- 💬 [Goa Slack Community](https://gophers.slack.com/channels/goa)
- 🐛 [Report Issues](https://github.com/goadesign/plugins/issues)
- 📦 [Plugin Source](https://github.com/goadesign/plugins/tree/main/testing)

---

*Built with ❤️ by the Goa community*