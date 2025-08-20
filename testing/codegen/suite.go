package codegen

import (
	"path/filepath"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

type (
	// suiteData contains data for generating test suite.
	suiteData struct {
		Package    string
		Service    *service.Data
		ImportPath string
		TestPkg    string // Package alias for the test harness
		NonStream  []*suiteMethodData
		Stream     []*suiteMethodData
		HasErrors  bool
		HasStreams bool
		HasHTTP    bool
		HasGRPC    bool
		HasJSONRPC bool
		UseTD      bool
		UseCtx     bool
	}

	// suiteTarget describes one transport/mode target to exercise for a method.
	suiteTarget struct {
		IsGRPC           bool
		IsHTTPPlain      bool
		IsHTTPServerSent bool
		IsHTTPWebSocket  bool
		IsJSONRPCPlain   bool
		IsJSONRPCSSE     bool
		IsJSONRPCWS      bool

		// Per-target stream kind flags (to support mixed transports)
		IsNoStream      bool
		IsClientStream  bool
		IsServerStream  bool
		IsBidirectional bool
	}

	// suiteMethodData describes a method in the test suite.
	suiteMethodData struct {
		// Method holds the authoritative Goa method data.
		Method  *service.MethodData
		Targets []*suiteTarget
	}
)

// generateSuiteTopLevel generates an example test suite for a service at the top-level package (outside gen).
func generateSuiteTopLevel(genpkg string, examplePkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	fname := codegen.SnakeCase(svc.Name) + "_suite_test.go"
	path := filepath.Join(codegen.Gendir, "..", fname)
	pkg := examplePkg
	data := buildSuiteData(genpkg, pkg, root, svc)
	specs := []*codegen.ImportSpec{
		{Path: "context"},
		{Path: "io"},
		{Path: "testing"},
		{Path: "time"},
		{Path: filepath.Join(genpkg, data.Service.PathName), Name: data.Service.PkgName},
		{Path: filepath.Join(genpkg, data.Service.PathName, data.Service.PathName+"test"), Name: data.Service.VarName+"test"},
	}
	sections := []*codegen.SectionTemplate{
		// Use empty title to generate header without "DO NOT EDIT" comment, like Goa examples
		codegen.Header("", pkg, specs),
		{
			Name:   "suite-test",
			Source: testingTemplates.Read(suiteTestT),
			FuncMap: template.FuncMap{
				"goify": codegen.Goify,
			},
			Data: data,
		},
	}
	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
		SkipExist:        true, // Don't overwrite existing test files
	}
}

// buildSuiteData builds the template data for test suite from DSL only.
func buildSuiteData(genpkg, pkg string, root *expr.RootExpr, svc *expr.ServiceExpr) *suiteData {
	sd := service.NewServicesData(root).Get(svc.Name)
	data := &suiteData{
		Package:    pkg,
		Service:    sd,
		ImportPath: filepath.Join(genpkg, sd.PathName),
		TestPkg:    sd.VarName + "test",
		NonStream:  make([]*suiteMethodData, 0, len(svc.Methods)),
		Stream:     make([]*suiteMethodData, 0, len(svc.Methods)),
		HasErrors:  hasErrors(svc),
		HasStreams: hasStreams(svc),
		HasHTTP:    hasHTTPTransport(root, svc),
		HasGRPC:    hasGRPCTransport(root, svc),
		HasJSONRPC: hasJSONRPCTransport(root, svc),
	}
	for _, m := range svc.Methods {
		md := &suiteMethodData{Method: sd.Method(m.Name)}
		// Compute targets for all transports the method supports
		var targets []*suiteTarget
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
		if hasMethodHTTP(root, svc, m) {
			if m.Stream == expr.NoStreamKind || m.Stream == 0 { // 0 is the zero value, treat as NoStreamKind
				// Plain HTTP, non-stream
				targets = append(targets, &suiteTarget{IsHTTPPlain: true, IsNoStream: true})
			} else {
				// Inspect HTTP endpoints for SSE; default to WebSocket when streaming
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
					// SSE is server-stream
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
		if hasMethodJSONRPC(root, svc, m) && md.Method != nil {
			if md.Method.IsJSONRPCSSE {
				// JSON-RPC SSE is modeled as server stream regardless of m.Stream
				targets = append(targets, &suiteTarget{IsJSONRPCSSE: true, IsServerStream: true})
			}
			if md.Method.IsJSONRPCWebSocket {
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
			if !md.Method.IsJSONRPCSSE && !md.Method.IsJSONRPCWebSocket {
				// Plain JSON-RPC non-stream
				targets = append(targets, &suiteTarget{IsJSONRPCPlain: true, IsNoStream: true})
			}
		}
		md.Targets = targets
		// Classify into NonStream and Stream based on targets (supports mixed transports)
		anyNonStream := false
		anyStream := false
		for _, t := range targets {
			if t.IsNoStream {
				anyNonStream = true
			}
			if t.IsClientStream || t.IsServerStream || t.IsBidirectional {
				anyStream = true
			}
		}
		if anyNonStream {
			data.NonStream = append(data.NonStream, md)
		}
		if anyStream {
			data.Stream = append(data.Stream, md)
		}

		// Determine if test data is needed (payloads passed to wrappers)
		if md.Method != nil {
			for _, t := range targets {
				if (t.IsNoStream && md.Method.PayloadEx != nil) || (t.IsServerStream && md.Method.PayloadEx != nil) {
					data.UseTD = true
					break
				}
			}
		}
	}
	data.UseCtx = len(svc.Methods) > 0

	return data
}
