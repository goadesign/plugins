package codegen

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

type (
	// clientData contains data for generating test client code.
	// It embeds Goa's service.Data to leverage existing codegen structures.
	clientData struct {
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
		// Methods with client-specific extensions
		Methods []*clientMethodData
	}

	// clientMethodData extends Goa's MethodData for client generation.
	clientMethodData struct {
		// Embed the base method data from Goa's codegen
		*service.MethodData
		// Targets lists all supported transport combinations
		Targets []*suiteTarget
		// Additional computed fields for client generation
		HasHTTP     bool // Method has HTTP transport (any variant)
		HasGRPC     bool // Method has gRPC transport
		HasJSONRPC  bool // Method has JSON-RPC transport (any variant)
		IsStreaming bool // Method has streaming
		// PkgResultRef is the package-qualified result type reference
		PkgResultRef string
	}
)

// generateClient generates the test client file for a service.
func generateClient(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	path := filepath.Join(testingPath(genpkg, svc), "client.go")

	if svcData == nil {
		return nil
	}

	data := buildClientData(svcData, root, svc)

	specs := []*codegen.ImportSpec{
		{Path: "context"},
		{Path: "fmt"},
		{Path: "testing"},
		{Path: "time"},
		{Path: filepath.Join(genpkg, codegen.SnakeCase(svc.Name)), Name: svcData.PkgName},
	}

	// Add transport-specific imports based on what the service uses
	if data.HasHTTP {
		specs = append(specs,
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "http", codegen.SnakeCase(svc.Name), "client"), Name: "httpcli"},
			&codegen.ImportSpec{Path: "goa.design/goa/v3/http", Name: "goahttp"},
		)
	}
	if data.HasGRPC {
		specs = append(specs,
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "grpc", codegen.SnakeCase(svc.Name), "client"), Name: "grpccli"},
			&codegen.ImportSpec{Path: "google.golang.org/grpc"},
		)
	}
	if data.HasJSONRPC {
		specs = append(specs,
			&codegen.ImportSpec{Path: filepath.Join(genpkg, "jsonrpc", codegen.SnakeCase(svc.Name), "client"), Name: "jsonrpccli"},
		)
	}

	sections := []*codegen.SectionTemplate{
		codegen.Header(fmt.Sprintf("Test client for %s service", svc.Name), codegen.SnakeCase(svc.Name)+"test", specs),
		{
			Name:   "client-struct",
			Source: testingTemplates.Read("client_struct"),
			Data:   data,
		},
		{
			Name:   "client-transport-selectors",
			Source: testingTemplates.Read("client_transport_selectors"),
			Data:   data,
		},
		{
			Name:   "client-stream-wrappers",
			Source: testingTemplates.Read("client_stream_wrappers"),
			Data:   data,
		},
		{
			Name:   "client-methods",
			Source: testingTemplates.Read("client_methods"),
			Data:   data,
		},
	}

	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	}
}

// buildClientData builds the template data for the client.
func buildClientData(svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) *clientData {
	// Create client data using Goa's service.Data directly
	data := &clientData{
		Data:        svcData,
		ServiceExpr: svc,
		HasHTTP:     hasHTTPTransport(root, svc),
		HasGRPC:     hasGRPCTransport(root, svc),
		HasJSONRPC:  hasJSONRPCTransport(root, svc),
		HasStreams:  hasStreams(svc),
		Methods:     make([]*clientMethodData, 0, len(svcData.Methods)),
	}

	// Create a name scope for type reference generation
	scope := codegen.NewNameScope()

	// Build method data with client-specific extensions
	for i, m := range svc.Methods {
		md := svcData.Methods[i]

		// Build targets for this method
		targets := buildMethodTargets(root, svc, m, md)

		// Create client method data
		cmd := &clientMethodData{
			MethodData:  md,
			Targets:     targets,
			IsStreaming: md.StreamKind != expr.NoStreamKind,
		}

		// Compute package-qualified result reference using Goa's GoFullTypeRef
		// This properly handles arrays, maps, primitives, and user types
		if m.Result != nil && m.Result.Type != expr.Empty {
			cmd.PkgResultRef = scope.GoFullTypeRef(m.Result, svcData.PkgName)
		}

		// Analyze targets to determine available transports
		for _, target := range targets {
			if target.IsGRPC {
				cmd.HasGRPC = true
			}
			if target.IsHTTPPlain || target.IsHTTPServerSent || target.IsHTTPWebSocket {
				cmd.HasHTTP = true
			}
			if target.IsJSONRPCPlain || target.IsJSONRPCSSE || target.IsJSONRPCWS {
				cmd.HasJSONRPC = true
			}
		}

		data.Methods = append(data.Methods, cmd)
	}

	return data
}

// buildMethodTargets builds the transport targets for a method.
// This is shared between harness and client generation.
func buildMethodTargets(root *expr.RootExpr, svc *expr.ServiceExpr, m *expr.MethodExpr, md *service.MethodData) []*suiteTarget {
	var targets []*suiteTarget

	// Check gRPC transport
	if hasMethodGRPC(root, svc, m) {
		g := &suiteTarget{IsGRPC: true}
		switch m.Stream {
		case expr.NoStreamKind, 0: // 0 is the zero value, treat as NoStreamKind
			g.IsNoStream = true
		case expr.ClientStreamKind:
			g.IsClientStream = true
		case expr.ServerStreamKind:
			g.IsServerStream = true
		case expr.BidirectionalStreamKind:
			g.IsBidirectional = true
		}
		targets = append(targets, g)
	}

	// Check HTTP transport
	if hasMethodHTTP(root, svc, m) {
		if m.Stream == expr.NoStreamKind || m.Stream == 0 { // 0 is the zero value, treat as NoStreamKind
			targets = append(targets, &suiteTarget{IsHTTPPlain: true, IsNoStream: true})
		} else {
			// Check for SSE
			isSSE := false
			if root.API != nil && root.API.HTTP != nil {
				for _, hs := range root.API.HTTP.Services {
					if hs.Name() != svc.Name {
						continue
					}
					if he := hs.Endpoint(m.Name); he != nil {
						isSSE = he.SSE != nil
						break
					}
				}
			}
			if isSSE {
				targets = append(targets, &suiteTarget{IsHTTPServerSent: true, IsServerStream: true})
			} else {
				ws := &suiteTarget{IsHTTPWebSocket: true}
				switch m.Stream {
				case expr.ClientStreamKind:
					ws.IsClientStream = true
				case expr.ServerStreamKind:
					ws.IsServerStream = true
				case expr.BidirectionalStreamKind:
					ws.IsBidirectional = true
				}
				targets = append(targets, ws)
			}
		}
	}

	// Check JSON-RPC transport
	if hasMethodJSONRPC(root, svc, m) && md != nil {
		switch {
		case md.IsJSONRPCSSE:
			targets = append(targets, &suiteTarget{IsJSONRPCSSE: true, IsServerStream: true})
		}
		if md.IsJSONRPCWebSocket {
			ws := &suiteTarget{IsJSONRPCWS: true}
			switch m.Stream {
			case expr.ClientStreamKind:
				ws.IsClientStream = true
			case expr.ServerStreamKind:
				ws.IsServerStream = true
			case expr.BidirectionalStreamKind:
				ws.IsBidirectional = true
			}
			targets = append(targets, ws)
		}
		if !md.IsJSONRPCSSE && !md.IsJSONRPCWebSocket {
			targets = append(targets, &suiteTarget{IsJSONRPCPlain: true, IsNoStream: true})
		}
	}

	return targets
}
