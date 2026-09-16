// This file renders the in-process servers and clients used by generated
// service tests for HTTP, gRPC, and JSON-RPC transports.
package codegen

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

type (
	// harnessData contains the data used to generate test harness code.
	// It embeds Goa's service.Data to leverage existing codegen structures.
	harnessData struct {
		// Embed the base service data from Goa's codegen
		*service.Data
		// Service expression for additional metadata
		ServiceExpr *expr.ServiceExpr
		// HasHTTP indicates if service has HTTP transport
		HasHTTP bool
		// HasGRPC indicates if service has gRPC transport
		HasGRPC bool
		// HasJSONRPC indicates if service has JSON-RPC transport
		HasJSONRPC bool
		// HasStreams indicates if service has streaming methods
		HasStreams bool
		// Methods with transport target information
		Methods []*harnessMethodData
	}

	// harnessMethodData extends Goa's MethodData with transport targets.
	// This allows us to leverage existing method data while adding test-specific info.
	harnessMethodData struct {
		// Embed the base method data from Goa's codegen
		*service.MethodData
		// Targets lists all supported transport combinations for this method
		Targets []*suiteTarget
	}
)

// generateHarness generates the test harness file for a service.
func generateHarness(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, transports *methodTransports) *codegen.File {
	path := filepath.Join(testingPath(svcData), "harness.go")

	if svcData == nil {
		return nil
	}

	data := buildHarnessData(svcData, root, svc, transports)

	specs := []*codegen.ImportSpec{
		{Path: "context"},
		{Path: "encoding/json"},
		{Path: "fmt"},
		{Path: "io"},
		{Path: "net/http"},
		{Path: "net/http/httptest"},
		{Path: "strings"},
		{Path: "testing"},
		{Path: "time"},
		{Path: filepath.Join(genpkg, svcData.PathName), Name: svcData.PkgName},
		{Path: "goa.design/goa/v3/pkg"},
	}

	// Add transport-specific imports
	if data.HasHTTP {
		specs = append(specs,
			&codegen.ImportSpec{Path: "bytes"},
			&codegen.ImportSpec{Path: "bufio"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "http", svcData.PathName, "server"), Name: "httpsvr"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "http", svcData.PathName, "client"), Name: "httpcli"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/http", Name: "goahttp"},
			&codegen.ImportSpec{Path: "net/url"},
			&codegen.ImportSpec{Path: "github.com/gorilla/websocket"},
		)
	}
	if data.HasGRPC {
		specs = append(specs,
			&codegen.ImportSpec{Path: "net"},
			&codegen.ImportSpec{Path: "google.golang.org/grpc"},
			&codegen.ImportSpec{Path: "google.golang.org/grpc/test/bufconn"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/grpc", Name: "goagrpc"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "grpc", svcData.PathName, "server"), Name: "grpcsvr"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "grpc", svcData.PathName, "client"), Name: "grpccli"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "grpc", svcData.PathName, "pb"), Name: svcData.PkgName + "pb"},
		)
	}
	if data.HasJSONRPC {
		specs = append(specs,
			&codegen.ImportSpec{Path: "net/url"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "jsonrpc", svcData.PathName, "server"), Name: "jsonrpcsvr"},
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "jsonrpc", svcData.PathName, "client"), Name: "jsonrpccli"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/jsonrpc", Name: "jsonrpc"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/http", Name: "goahttp"},
		)
	}

	sections := []*codegen.SectionTemplate{
		codegen.Header(fmt.Sprintf("Test harness for %s service", svc.Name), svcData.PathName+"test", specs),
		{
			Name:   "harness-struct",
			Source: testingTemplates.Read(harnessStructT),
			Data:   data,
		},
		{
			Name:   "harness-constructor",
			Source: testingTemplates.Read(harnessConstructorT),
			Data:   data,
		},
	}

	// Transport setup sections
	if data.HasHTTP {
		sections = append(sections, &codegen.SectionTemplate{Name: "http-harness", Source: testingTemplates.Read(httpHarnessT), Data: data})
	}
	if data.HasGRPC {
		sections = append(sections, &codegen.SectionTemplate{Name: "grpc-harness", Source: testingTemplates.Read(grpcHarnessT), Data: data})
	}
	if data.HasJSONRPC {
		sections = append(sections, &codegen.SectionTemplate{Name: "jsonrpc-harness", Source: testingTemplates.Read(jsonrpcHarnessT), Data: data})
	}

	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	}
}

// buildHarnessData builds the template data for the harness.
func buildHarnessData(svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr, transports *methodTransports) *harnessData {
	// Create harness data using Goa's service.Data directly
	data := &harnessData{
		Data:        svcData,
		ServiceExpr: svc,
		HasHTTP:     hasHTTPTransport(root, svc),
		HasGRPC:     hasGRPCTransport(root, svc),
		HasJSONRPC:  hasJSONRPCTransport(root, svc),
		HasStreams:  hasStreams(svc),
		Methods:     make([]*harnessMethodData, 0, len(svcData.Methods)),
	}

	// Build method data with test-specific extensions
	for i, m := range svc.Methods {
		md := svcData.Methods[i]

		// Build targets for this method using shared function
		targets := buildMethodTargets(root, svc, m, transports)

		// Create harness method data
		hmd := &harnessMethodData{
			MethodData: md,
			Targets:    targets,
		}

		data.Methods = append(data.Methods, hmd)
	}

	return data
}

// hasHTTPTransport checks if the service has HTTP transport.
func hasHTTPTransport(root *expr.RootExpr, svc *expr.ServiceExpr) bool {
	// Check if service is in HTTP services
	if root != nil && root.API != nil && root.API.HTTP != nil {
		for _, s := range root.API.HTTP.Services {
			if s.Name() == svc.Name {
				return true
			}
		}
	}
	return false
}

// hasGRPCTransport checks if the service has gRPC transport.
func hasGRPCTransport(root *expr.RootExpr, svc *expr.ServiceExpr) bool {
	// Check if service is in gRPC services
	if root != nil && root.API != nil && root.API.GRPC != nil {
		for _, s := range root.API.GRPC.Services {
			if s.Name() == svc.Name {
				return true
			}
		}
	}
	return false
}

// hasJSONRPCTransport checks if the service has JSON-RPC transport.
func hasJSONRPCTransport(root *expr.RootExpr, svc *expr.ServiceExpr) bool {
	if root != nil && root.API != nil && root.API.JSONRPC != nil {
		for _, s := range root.API.JSONRPC.Services {
			if s.Name() == svc.Name {
				return true
			}
		}
	}
	return false
}

// testingPath returns the output path next to the generated service package.
func testingPath(svc *service.Data) string {
	return filepath.Join(codegen.Gendir, svc.PathName, svc.PathName+"test")
}

// hasStreams checks if the service has streaming methods.
func hasStreams(svc *expr.ServiceExpr) bool {
	for _, m := range svc.Methods {
		if m.IsStreaming() {
			return true
		}
	}
	return false
}

// hasMethodHTTP checks whether the given method is bound to HTTP transport.
func hasMethodHTTP(root *expr.RootExpr, svc *expr.ServiceExpr, m *expr.MethodExpr) bool {
	if root != nil && root.API != nil && root.API.HTTP != nil {
		for _, hs := range root.API.HTTP.Services {
			if hs.Name() != svc.Name {
				continue
			}
			for _, hm := range hs.HTTPEndpoints {
				if hm.MethodExpr == m {
					return true
				}
			}
		}
	}
	return false
}

// hasMethodGRPC checks whether the given method is bound to gRPC transport.
func hasMethodGRPC(root *expr.RootExpr, svc *expr.ServiceExpr, m *expr.MethodExpr) bool {
	if root != nil && root.API != nil && root.API.GRPC != nil {
		for _, gs := range root.API.GRPC.Services {
			if gs.Name() != svc.Name {
				continue
			}
			for _, gm := range gs.GRPCEndpoints {
				if gm.MethodExpr == m {
					return true
				}
			}
		}
	}
	return false
}

// hasPayloads checks if the service has methods with payloads.
func hasPayloads(svc *expr.ServiceExpr) bool {
	for _, m := range svc.Methods {
		if m.Payload.Type != expr.Empty || m.StreamingPayload.Type != expr.Empty {
			return true
		}
	}
	return false
}
