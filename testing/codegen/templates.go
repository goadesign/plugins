package codegen

import (
	"embed"

	"goa.design/goa/v3/codegen/template"
)

// Template constants
const (
	// Core harness templates
	harnessStructT      = "harness_struct"
	harnessConstructorT = "harness_constructor"
	implementHookT      = "implement_hook"

	// Transport-specific harness templates
	httpHarnessT    = "http_harness"
	grpcHarnessT    = "grpc_harness"
	jsonrpcHarnessT = "jsonrpc_harness"

	// Error testing templates
	errorHelpersT = "error_helpers"

	// Test data templates
	testdataGeneratorsT = "testdata_generators"

	// Test suite templates
	suiteTestT = "suite_test"
)

//go:embed templates/*
var templateFS embed.FS

// testingTemplates provides access to the testing package templates.
var testingTemplates = &template.TemplateReader{FS: templateFS}
