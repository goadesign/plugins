# Calculator Service Example - Goa Testing Plugin Showcase

This example demonstrates all the features of the Goa testing plugin through a simple calculator service. It's designed to be educational rather than realistic, showcasing how to use the testing plugin effectively.

## 📚 What This Example Demonstrates

### 1. **Basic Testing**
- Simple test scenarios with expected results
- Multi-step test workflows
- Edge case testing

### 2. **Transport Selection**
- Testing with HTTP and gRPC transports
- Transport-specific scenarios
- Automatic transport selection

### 3. **Error Handling**
- Expected error scenarios (division by zero)
- Error message validation
- DSL validation errors

### 4. **Custom Validators**
- Business logic validation beyond DSL constraints
- Precision checking for floating-point operations
- Performance validation (timing checks)
- Statistical accuracy verification
- Outlier detection

### 5. **Timeout Configuration**
- Scenario-level timeouts (apply to all steps)
- Step-level timeouts (override scenario timeout)
- Context cancellation handling

### 6. **Streaming Operations**
- Bidirectional streaming with `batch_add`
- Stream validation patterns

## 🚀 Quick Start

### 1. Generate the code
```bash
goa gen goa.design/plugins/v3/testing/examples/calculator/design
goa example goa.design/plugins/v3/testing/examples/calculator/design
```

### 2. Run the tests
```bash
go test -v
```

### 3. Run specific scenarios
```bash
go test -v -run TestScenarios/simple_addition
go test -v -run TestScenarios/division_with_validation
```

## 📁 Project Structure

```
calculator/
├── design/
│   └── design.go           # Service DSL definition
├── gen/                    # Generated code
│   └── calculator/
│       └── calculatortest/ # Testing plugin generated code
│           ├── scenarios.go     # Scenario runner
│           ├── client.go        # Test client
│           ├── harness.go       # Test harness
│           └── validators_test.go # Custom validators
├── scenarios.yaml          # Test scenarios
├── calculator.go           # Service implementation
└── calculator_suite_test.go # Test suite
```

## 🎯 Key Features Explained

### Scenario Configuration (`scenarios.yaml`)

The YAML file defines all test scenarios declaratively:

```yaml
scenarios:
  - name: "simple_addition"
    steps:
      - method: add
        payload:
          a: 5
          b: 3
        expect:
          result:
            result: 8
            operation: "5 + 3"
```

### Custom Validators

Validators allow business logic validation beyond what the DSL provides:

```go
// ValidatePrecision checks division precision
func ValidatePrecision(result *calculator.DivideResult, expected map[string]any) error {
    expectedValue := 100.0 / 3.0
    tolerance := 0.0001
    
    if math.Abs(result.Result - expectedValue) > tolerance {
        return fmt.Errorf("precision error: expected ~%f, got %f", expectedValue, result.Result)
    }
    return nil
}
```

Use in scenarios.yaml:
```yaml
- method: divide
  payload:
    dividend: 100
    divisor: 3
  expect:
    validator: ValidatePrecision
```

### Timeout Configuration

Control test execution time:

```yaml
scenarios:
  - name: "fast_operations"
    timeout: "1s"  # Scenario-level default
    steps:
      - method: factorial
        payload:
          n: 15
        timeout: "500ms"  # Step-level override
```

### Transport Selection

Force specific transports for testing:

```yaml
scenarios:
  - name: "grpc_test"
    transport: grpc  # Use gRPC for all steps
    steps:
      - method: add
        # ...
```

## 📊 Available Test Scenarios

| Scenario | Description | Features Demonstrated |
|----------|-------------|----------------------|
| `simple_addition` | Basic addition test | Basic validation |
| `addition_via_grpc` | Addition using gRPC | Transport selection |
| `division_by_zero_error` | Error handling | Expected errors |
| `division_with_validation` | Custom precision check | Custom validators |
| `fast_factorial` | Small factorial | Timeout handling |
| `large_factorial` | Larger factorial with timing | Performance validation |
| `statistics_validation` | Statistical calculations | Complex validation |
| `edge_cases` | Various edge cases | Comprehensive testing |
| `calculation_workflow` | Multi-step calculations | Workflow testing |
| `statistics_outliers` | Outlier detection | Advanced validation |
| `floating_point_precision` | Float precision issues | Precision handling |

## 🧪 Running Tests

### Run all tests
```bash
go test -v
```

### Run specific scenario
```bash
go test -v -run TestScenarios/division_with_validation
```

### Run with coverage
```bash
go test -v -cover
```

### Run with timeout
```bash
go test -v -timeout 30s
```

## 🔧 Customization

### Adding New Validators

1. Add validator function to `validators_test.go`:
```go
func ValidateMyLogic(result *calculator.SomeResult, expected map[string]any) error {
    // Your validation logic
    return nil
}
```

2. Reference in `scenarios.yaml`:
```yaml
expect:
  validator: ValidateMyLogic
```

### Adding New Scenarios

Add to `scenarios.yaml`:
```yaml
scenarios:
  - name: "my_new_test"
    description: "Test description"
    timeout: "2s"
    steps:
      - method: add
        payload:
          a: 1
          b: 2
        expect:
          result:
            result: 3
```

## 💡 Tips

1. **Use validators for business logic** - DSL handles type validation, validators handle business rules
2. **Set appropriate timeouts** - Prevent hanging tests with reasonable timeouts
3. **Test error paths** - Don't just test happy paths
4. **Group related tests** - Use multi-step scenarios for workflows
5. **Document validators** - Explain what each validator checks

## 🐛 Troubleshooting

### Tests fail with timeout
- Increase timeout in scenarios.yaml
- Check service implementation for infinite loops
- Verify context cancellation is handled

### Validator not found
- Ensure validator is exported (starts with capital letter)
- Check validator is in `_test.go` file in the test package
- Verify validator name matches exactly in YAML

### Transport errors
- Ensure service implements both HTTP and gRPC
- Check transport name is valid (http, grpc, auto)
- Verify method supports requested transport

## 📚 Learn More

- [Goa Testing Plugin Documentation](../../README.md)
- [Goa Framework](https://goa.design)
- [Writing Effective Tests](https://goa.design/learn/testing)