// This file converts Goa's finalized JSON-RPC plan into the small set of
// transport choices used by testing templates. The templates receive only the
// branches that the design can use, so generated tests contain no runtime
// transport guessing.
package codegen

import (
	"fmt"

	"goa.design/goa/v3/expr"
	jsonrpccodegen "goa.design/goa/v3/jsonrpc/codegen"
)

type (
	// methodTransports records the JSON-RPC transport selected for each method.
	// HTTP and gRPC choices remain available directly from their design entries.
	methodTransports struct {
		jsonRPC map[*expr.MethodExpr]jsonRPCMethodTransport
	}

	// jsonRPCMethodTransport records whether one JSON-RPC method sends a
	// server-sent event stream. A false value means an ordinary JSON-RPC call.
	jsonRPCMethodTransport struct {
		serverSentEvents bool
	}
)

// plannedMethodTransports reads each finalized JSON-RPC endpoint and returns
// the transport choice used to specialize the generated testing helpers.
func plannedMethodTransports(root *expr.RootExpr, plan *jsonrpccodegen.Plan) *methodTransports {
	transports := &methodTransports{jsonRPC: make(map[*expr.MethodExpr]jsonRPCMethodTransport)}
	if root.API.JSONRPC == nil {
		return transports
	}
	for _, transportService := range root.API.JSONRPC.Services {
		service, ok := plan.Service(transportService)
		if !ok {
			panic(fmt.Sprintf("JSON-RPC plan is missing service %q", transportService.Name()))
		}
		serviceExpr := root.Service(transportService.Name())
		for _, endpoint := range service.Endpoints {
			method := serviceExpr.Method(endpoint.Method.Name)
			transports.jsonRPC[method] = jsonRPCMethodTransport{serverSentEvents: endpoint.SSE != nil}
		}
	}
	return transports
}

// designedMethodTransports reads the prepared JSON-RPC design for callers of
// the released callback API, which does not receive Goa's generation plan.
func designedMethodTransports(root *expr.RootExpr) *methodTransports {
	transports := &methodTransports{jsonRPC: make(map[*expr.MethodExpr]jsonRPCMethodTransport)}
	if root.API.JSONRPC == nil {
		return transports
	}
	for _, service := range root.API.JSONRPC.Services {
		for _, endpoint := range service.HTTPEndpoints {
			transports.jsonRPC[endpoint.MethodExpr] = jsonRPCMethodTransport{serverSentEvents: endpoint.SSE != nil}
		}
	}
	return transports
}
