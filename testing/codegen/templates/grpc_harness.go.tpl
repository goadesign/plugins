{{ printf "setupGRPC initializes the gRPC test server and client connection." | comment }}
func (h *Harness) setupGRPC() {
	// Create buffer for in-memory connection
	lis := bufconn.Listen(1024 * 1024)
	
	// Create gRPC server
	h.grpcSvr = grpc.NewServer()
	
	// Register service
	endpoints := {{ .PkgName }}.NewEndpoints(h.service)
	server := grpcsvr.New(endpoints, nil, nil)
	{{ .PkgName }}pb.Register{{ .StructName }}Server(h.grpcSvr, server)
	
	// Start server in background
	go func() {
		if err := h.grpcSvr.Serve(lis); err != nil {
			fmt.Printf("gRPC server stopped (%v)\n", err)
		}
	}()
	
	// Create client connection
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		h.t.Fatalf("Failed to create gRPC client connection: %v", err)
	}
	
	h.grpcConn = conn
}

{{ printf "GRPCConn returns the gRPC client connection." | comment }}
func (h *Harness) GRPCConn() *grpc.ClientConn {
	if h.grpcConn == nil {
		h.t.Fatal("gRPC transport not configured")
	}
	return h.grpcConn
}

{{ printf "getGRPCClientImpl returns the underlying gRPC client implementation." | comment }}
func (h *Harness) getGRPCClientImpl() *grpccli.Client {
	conn := h.GRPCConn()
	return grpccli.NewClient(conn)
}

{{ printf "GRPCClient creates gRPC client endpoints for the service." | comment }}
func (h *Harness) GRPCClient() *{{ .PkgName }}.Endpoints {
	client := h.getGRPCClientImpl()
	return &{{ .PkgName }}.Endpoints{
		{{- range .Methods }}
		{{- $method := . }}
		{{- range .Targets }}
		{{- if .IsGRPC }}
		{{ $method.VarName }}: client.{{ $method.VarName }}(),
		{{- break }}
		{{- end }}
		{{- end }}
		{{- end }}
	}
}
