package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"sync"

	"goa.design/clue/debug"
	"goa.design/clue/log"
	test_http_grpcpb "goa.design/plugins/v3/testing/examples/httpgrpc/gen/grpc/test_http_grpc/pb"
	testhttpgrpcsvr "goa.design/plugins/v3/testing/examples/httpgrpc/gen/grpc/test_http_grpc/server"
	testhttpgrpc "goa.design/plugins/v3/testing/examples/httpgrpc/gen/test_http_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// handleGRPCServer starts configures and starts a gRPC server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleGRPCServer(ctx context.Context, u *url.URL, testHTTPGrpcEndpoints *testhttpgrpc.Endpoints, wg *sync.WaitGroup, errc chan error, dbg bool) {

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to gRPC requests and
	// responses.
	var (
		testHTTPGrpcServer *testhttpgrpcsvr.Server
	)
	{
		testHTTPGrpcServer = testhttpgrpcsvr.New(testHTTPGrpcEndpoints, nil, nil)
	}

	// Create interceptor which sets up the logger in each request context.
	chain := grpc.ChainUnaryInterceptor(log.UnaryServerInterceptor(ctx))
	if dbg {
		// Log request and response content if debug logs are enabled.
		chain = grpc.ChainUnaryInterceptor(log.UnaryServerInterceptor(ctx), debug.UnaryServerInterceptor())
	}
	streamchain := grpc.ChainStreamInterceptor(log.StreamServerInterceptor(ctx))
	if dbg {
		streamchain = grpc.ChainStreamInterceptor(log.StreamServerInterceptor(ctx), debug.StreamServerInterceptor())
	}

	// Initialize gRPC server
	srv := grpc.NewServer(chain, streamchain)

	// Register the servers.
	test_http_grpcpb.RegisterTestHTTPGrpcServer(srv, testHTTPGrpcServer)

	for svc, info := range srv.GetServiceInfo() {
		for _, m := range info.Methods {
			log.Printf(ctx, "serving gRPC method %s", svc+"/"+m.Name)
		}
	}

	// Register the server reflection service on the server.
	// See https://grpc.github.io/grpc/core/md_doc_server-reflection.html.
	reflection.Register(srv)

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start gRPC server in a separate goroutine.
		go func() {
			lis, err := net.Listen("tcp", u.Host)
			if err != nil {
				errc <- err
			}
			if lis == nil {
				errc <- fmt.Errorf("failed to listen on %q", u.Host)
			}
			log.Printf(ctx, "gRPC server listening on %q", u.Host)
			errc <- srv.Serve(lis)
		}()

		<-ctx.Done()
		log.Printf(ctx, "shutting down gRPC server at %q", u.Host)
		srv.Stop()
	}()
}
