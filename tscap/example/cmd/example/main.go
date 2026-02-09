package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/tscap/example"
	tscapsvc "goa.design/plugins/v3/tscap/example/gen/tscap"
	tscapsvr "goa.design/plugins/v3/tscap/example/gen/http/tscap/server"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "HTTP listen address")
		debug = flag.Bool("debug", false, "print all request headers")
	)
	flag.Parse()

	logger := log.New(os.Stderr, "[tscap] ", log.Ltime)

	// Create service
	svc := example.NewTscap(logger)
	endpoints := tscapsvc.NewEndpoints(svc)

	// Create transport
	mux := goahttp.NewMuxer()
	dec := goahttp.RequestDecoder
	enc := goahttp.ResponseEncoder
	svr := tscapsvr.New(endpoints, mux, dec, enc, nil, nil)
	tscapsvr.Mount(mux, svr)

	// Wrap mux with logging middleware
	handler := loggingMiddleware(logger, *debug, mux)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         *addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server
	go func() {
		logger.Printf("listening on %s", *addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Print("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal(err)
	}
	fmt.Println("done")
}

func loggingMiddleware(logger *log.Logger, debug bool, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if debug {
			logger.Printf("%s %s headers:", r.Method, r.URL.Path)
			for name, values := range r.Header {
				for _, v := range values {
					logger.Printf("  %s: %s", name, v)
				}
			}
		} else {
			caps := r.Header.Get("Tailscale-App-Capabilities")
			if caps == "" {
				logger.Printf("%s %s - no caps header", r.Method, r.URL.Path)
			} else {
				preview := caps
				if len(preview) > 100 {
					preview = preview[:100] + "..."
				}
				logger.Printf("%s %s - caps: %s", r.Method, r.URL.Path, preview)
			}
		}
		next.ServeHTTP(w, r)
	})
}
