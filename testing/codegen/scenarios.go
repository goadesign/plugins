// This file renders scenario runners and example scenario files from finalized
// Goa service data and transport choices.
package codegen

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
	"gopkg.in/yaml.v3"
)

// generateScenarios generates the scenario runner for a service.
func generateScenarios(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, transports *methodTransports) []*codegen.File {
	if svcData == nil {
		return nil
	}

	data := buildScenariosData(svcData, root, svc, transports)
	files := make([]*codegen.File, 0, 2)

	// Generate main scenarios runner file
	path := filepath.Join(testingPath(svcData), "scenarios.go")
	specs := []*codegen.ImportSpec{
		{Path: "context"},
		{Path: "encoding/json"},
		{Path: "fmt"},
		{Path: "os"},
		{Path: "strings"},
		{Path: "testing"},
		{Path: "time"},
		{Path: "gopkg.in/yaml.v3", Name: "yaml"},
		{Path: filepath.Join(genpkg, svcData.PathName), Name: svcData.PkgName},
	}

	// Add validator package import if specified in YAML
	// Only add if it's different from the current package
	currentPkgPath := filepath.Join(genpkg, svcData.PathName, svcData.PathName+"test")
	validatorInfo := ExtractValidatorsFromYAML()
	if validatorInfo.Path != "" && validatorInfo.Package != "" && validatorInfo.Path != currentPkgPath {
		specs = append(specs, &codegen.ImportSpec{
			Path: validatorInfo.Path,
			Name: validatorInfo.Package,
		})
		// Set the package name for use in template
		data.ValidatorPkg = validatorInfo.Package
	}

	sections := []*codegen.SectionTemplate{
		codegen.Header(fmt.Sprintf("Scenario runner for %s service", svc.Name), svcData.PathName+"test", specs),
		{
			Name:   "scenario-types",
			Source: testingTemplates.Read("scenario_types"),
			Data:   data,
		},
		{
			Name:   "scenario-runner",
			Source: testingTemplates.Read("scenario_runner"),
			Data:   data,
		},
	}

	files = append(files, &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	})

	return files
}

type (
	// scenariosData contains the data used to generate scenario runner code.
	// It embeds Goa's service.Data to leverage existing codegen structures.
	scenariosData struct {
		// Embed the base service data from Goa's codegen
		*service.Data
		// Service expression for additional metadata
		ServiceExpr *expr.ServiceExpr
		// Methods with scenario-specific extensions
		Methods []*scenarioMethodData
		// HasHTTP indicates if service has HTTP transport
		HasHTTP bool
		// HasGRPC indicates if service has gRPC transport
		HasGRPC bool
		// HasJSONRPC indicates if service has JSON-RPC transport
		HasJSONRPC bool
		// ValidTransports lists the transports used by at least one method.
		ValidTransports []*transportData
		// Validators found in scenarios.yaml
		Validators map[string][]string // method -> validator names
		// ValidatorPkg is the package name for validators
		ValidatorPkg string
		// ValidatorPath is the import path for the validator package
		ValidatorPath string
	}

	// scenarioMethodData extends Goa's MethodData for scenario generation.
	scenarioMethodData struct {
		// Embed the base method data from Goa's codegen
		*service.MethodData
		// Transports lists valid transport strings for YAML
		Transports []string
		// ResultTypeRef is the fully qualified Go reference to the result type
		// as seen from the generated test package (e.g. "*svc.Foo", "[]*svc.Bar").
		// This is used in generated type assertions for custom validators.
		ResultTypeRef string
		// IsClientStream reports that the client sends a sequence of payloads.
		IsClientStream bool
		// IsServerStream reports that the server sends a sequence of results.
		IsServerStream bool
		// IsBidirectional reports that both client and server send sequences.
		IsBidirectional bool
	}

	// transportData describes one scenario transport in generated help text.
	transportData struct {
		Name        string
		Description string
	}
)

// generateExampleScenarios generates an example scenarios.yaml file for a service.
func generateExampleScenarios(_ string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, transports *methodTransports) *codegen.File {
	path := filepath.Join(codegen.Gendir, "..", "scenarios.yaml")
	if svcData == nil {
		return nil
	}

	data := buildScenariosData(svcData, root, svc, transports)

	// For YAML files, we need to read the template directly since it's not a .go.tpl file
	tmplContent, err := templateFS.ReadFile("templates/example_scenarios.yaml.tpl")
	if err != nil {
		panic(fmt.Sprintf("failed to read example_scenarios.yaml.tpl: %v", err))
	}

	sections := []*codegen.SectionTemplate{
		{
			Name:   "example-scenarios",
			Source: string(tmplContent),
			Data:   data,
		},
	}

	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
		SkipExist:        true, // Don't overwrite existing scenarios file
	}
}

func buildScenariosData(svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, transports *methodTransports) *scenariosData {
	// Extract validator info from YAML
	validatorInfo := ExtractValidatorsFromYAML()

	data := &scenariosData{
		Data:            svcData,
		ServiceExpr:     svc,
		Methods:         make([]*scenarioMethodData, 0),
		HasHTTP:         hasHTTPTransport(root, svc),
		HasGRPC:         hasGRPCTransport(root, svc),
		HasJSONRPC:      hasJSONRPCTransport(root, svc),
		ValidTransports: make([]*transportData, 0),
		Validators:      validatorInfo.Validators,
		ValidatorPkg:    "", // Will be set below if it's a different package
		ValidatorPath:   validatorInfo.Path,
	}

	// Collect only transports that at least one generated method can use.
	allTransports := map[string]bool{"auto": true}

	// Build method data with available transports
	for i, m := range svc.Methods {
		methodData := svcData.Methods[i]

		// Build targets for this method using shared function
		targets := buildMethodTargets(root, svc, m, transports)

		md := &scenarioMethodData{
			MethodData: methodData,
			Transports: make([]string, 0),
		}
		switch m.Stream {
		case expr.ClientStreamKind:
			md.IsClientStream = true
		case expr.ServerStreamKind:
			md.IsServerStream = true
		case expr.BidirectionalStreamKind:
			md.IsBidirectional = true
		}
		// Compute fully qualified result type reference for type assertions.
		// This properly handles composite types like ArrayOf(...) without producing
		// invalid Go like "svc.[]T" (see issue #234).
		if m.Result != nil && m.Result.Type != expr.Empty {
			md.ResultTypeRef = svcData.Scope.GoFullTypeRef(m.Result, svcData.PkgName)
		}

		// Build list of valid transport strings based on targets
		methodTransports := make(map[string]bool)
		for _, target := range targets {
			switch {
			case target.IsGRPC:
				methodTransports["grpc"] = true
			case target.IsHTTPPlain:
				methodTransports["http"] = true
			case target.IsHTTPServerSent:
				methodTransports["http-sse"] = true
			case target.IsHTTPWebSocket:
				methodTransports["http-ws"] = true
			case target.IsJSONRPCPlain:
				methodTransports["jsonrpc"] = true
			case target.IsJSONRPCSSE:
				methodTransports["jsonrpc-sse"] = true
			}
		}

		// Convert set to sorted list
		md.Transports = slices.Sorted(maps.Keys(methodTransports))
		for transport := range methodTransports {
			allTransports[transport] = true
		}

		data.Methods = append(data.Methods, md)
	}
	for _, transport := range slices.Sorted(maps.Keys(allTransports)) {
		data.ValidTransports = append(data.ValidTransports, &transportData{
			Name:        transport,
			Description: transportDescription(transport),
		})
	}

	return data
}

// transportDescription explains one generated scenario transport.
func transportDescription(transport string) string {
	switch transport {
	case "auto":
		return "Use the first available transport"
	case "grpc":
		return "gRPC"
	case "http":
		return "HTTP"
	case "http-sse":
		return "HTTP server-sent events"
	case "http-ws":
		return "HTTP WebSocket"
	case "jsonrpc":
		return "JSON-RPC over HTTP"
	case "jsonrpc-sse":
		return "JSON-RPC over server-sent events"
	default:
		panic(fmt.Sprintf("testing plugin has no description for transport %q", transport))
	}
}

// ValidatorInfo holds validator configuration
type ValidatorInfo struct {
	Validators map[string][]string // method -> validator names
	Package    string              // package name containing validators
	Path       string              // import path for the package
}

// ExtractValidatorsFromYAML reads scenarios.yaml if it exists and extracts validator names and package
func ExtractValidatorsFromYAML() ValidatorInfo {
	info := ValidatorInfo{
		Validators: make(map[string][]string),
		Package:    "", // Empty means use current package
	}

	// Try to read scenarios.yaml from current directory
	data, err := os.ReadFile("scenarios.yaml")
	if err != nil {
		// File doesn't exist or can't be read, that's OK
		return info
	}

	if len(data) == 0 {
		// Empty file
		return info
	}

	var config struct {
		Validators struct {
			Package string `yaml:"package"`
			Path    string `yaml:"path"`
		} `yaml:"validators"`
		Scenarios []struct {
			Steps []struct {
				Method string `yaml:"method"`
				Expect struct {
					Validator string `yaml:"validator"`
				} `yaml:"expect"`
			} `yaml:"steps"`
		} `yaml:"scenarios"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		// Invalid YAML, skip
		return info
	}

	// Extract package info
	info.Package = config.Validators.Package
	info.Path = config.Validators.Path

	// Extract unique validators per method
	for _, scenario := range config.Scenarios {
		if scenario.Steps == nil {
			continue
		}
		for _, step := range scenario.Steps {
			if step.Method == "" || step.Expect.Validator == "" {
				continue
			}
			// Add to the list if not already there
			found := false
			for _, v := range info.Validators[step.Method] {
				if v == step.Expect.Validator {
					found = true
					break
				}
			}
			if !found {
				info.Validators[step.Method] = append(info.Validators[step.Method], step.Expect.Validator)
			}
		}
	}

	return info
}
