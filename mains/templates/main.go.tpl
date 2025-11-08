func main() {
    var (
        {{- if .HasHTTP }}
        httpaddr    = flag.String("http-addr", ":8080", "HTTP listen address")
        {{- end }}
        {{- if .HasGRPC }}
        grpcaddr    = flag.String("grpc-addr", ":9090", "gRPC listen address")
        {{- end }}
        metricsAddr = flag.String("metrics-addr", ":8081", "metrics listen address")
        coladdr     = flag.String("otel-addr", ":4317", "OpenTelemetry collector listen address")
        debugf      = flag.Bool("debug", false, "Enable debug logs")
    )
    flag.Parse()

    // 1. Create logger
    format := log.FormatJSON
    if log.IsTerminal() {
        format = log.FormatTerminal
    }
    ctx := log.Context(context.Background(), log.WithFormat(format), log.WithFunc(log.Span))
    {{- if eq .ServiceCount 1 }}
    ctx = log.With(ctx, log.KV{K: "svc", V: {{ (index .Services 0).GenPkg }}.ServiceName})
    {{- else }}
    ctx = log.With(ctx, log.KV{K: "svc", V: "{{ .ServerLabel }}"})
    {{- end }}
    if *debugf {
        ctx = log.Context(ctx, log.WithDebug())
        log.Debugf(ctx, "debug logs enabled")
    }

    // 2. Setup instrumentation (OTLP/GRPC)
    spanExporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint(*coladdr),
        otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf(ctx, err, "failed to initialize tracing")
    }
    defer func() {
        // Create new context in case the parent context has been canceled.
        ctx := log.Context(context.Background(), log.WithFormat(format))
        if err := spanExporter.Shutdown(ctx); err != nil {
            log.Errorf(ctx, err, "failed to shutdown tracing")
        }
    }()
    metricExporter, err := otlpmetricgrpc.New(ctx,
        otlpmetricgrpc.WithEndpoint(*coladdr),
        otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf(ctx, err, "failed to initialize metrics")
    }
    defer func() {
        // Create new context in case the parent context has been canceled.
        ctx := log.Context(context.Background(), log.WithFormat(format))
        if err := metricExporter.Shutdown(ctx); err != nil {
            log.Errorf(ctx, err, "failed to shutdown metrics")
        }
    }()

    {{- if eq .ServiceCount 1 }}
    cfg, err := clue.NewConfig(ctx,
        {{ (index .Services 0).GenPkg }}.ServiceName,
        {{ (index .Services 0).GenPkg }}.APIVersion,
        metricExporter,
        spanExporter,
    )
    {{- else }}
    cfg, err := clue.NewConfig(ctx,
        "{{ .ServerLabel }}",
        "v1",
        metricExporter,
        spanExporter,
    )
    {{- end }}
    if err != nil {
        log.Fatalf(ctx, err, "failed to initialize instrumentation")
    }
    clue.ConfigureOpenTelemetry(ctx, cfg)

    // 3. Create transport-agnostic HTTP client for downstream calls (instrumented)
    _ = &http.Client{ // example HTTP client; replace per-service as needed
        Transport: log.Client(
            otelhttp.NewTransport(
                http.DefaultTransport,
                otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
                    return otelhttptrace.NewClientTrace(ctx)
                }),
            ))}

    // 4. Mount health check & metrics on separate HTTP server
    check := health.Handler(health.NewChecker())
    check = log.HTTP(ctx)(check).(http.HandlerFunc) // Log health-check errors
    http.Handle("/healthz", check)
    http.Handle("/livez", check)
    metricsServer := &http.Server{Addr: *metricsAddr}

    // 5. Create services & endpoints
    {{- range .Services }}
    // {{ .Name }} service
    {{ .SvcVar }} := {{ $.APIPkg }}.New{{ .StructName }}()
    {{ .EpVar }} := {{ .GenPkg }}.NewEndpoints({{ .SvcVar }})
    {{ .EpVar }}.Use(debug.LogPayloads())
    {{ .EpVar }}.Use(log.Endpoint)
    {{- end }}

    {{- if .HasHTTP }}
    // 6. Create HTTP transport
    mux := goahttp.NewMuxer()
    debug.MountDebugLogEnabler(debug.Adapt(mux))
    debug.MountPprofHandlers(debug.Adapt(mux))
    handler := otelhttp.NewHandler(mux, {{ if eq .ServiceCount 1 }} {{ (index .Services 0).GenPkg }}.ServiceName {{ else }} "{{ .ServerLabel }}" {{ end }})
    handler = debug.HTTP()(handler) // Add debug endpoints
    handler = log.HTTP(ctx)(handler) // Add logger to request context

    {{- if .HasAnyWebSocket }}
    upgrader := &websocket.Upgrader{}
    {{- end }}

    {{- range .Services }}
    {{- if .HasHTTP }}
    // {{ .Name }} HTTP server
    {{- if .HasWebSocket }}
    {{ .SrvVar }} := {{ .GenHTTPPkg }}.New({{ .EpVar }}, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil, upgrader, nil{{- range .FileServerNils }}, nil{{- end }})
    {{- else }}
    {{ .SrvVar }} := {{ .GenHTTPPkg }}.New({{ .EpVar }}, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil{{- range .FileServerNils }}, nil{{- end }})
    {{- end }}
    {{ .GenHTTPPkg }}.Mount(mux, {{ .SrvVar }})
    for _, m := range {{ .SrvVar }}.Mounts {
        log.Print(ctx, log.KV{K: "method", V: m.Method}, log.KV{K: "endpoint", V: m.Verb + " " + m.Pattern})
    }
    {{- end }}
    {{- end }}

    httpServer := &http.Server{Addr: *httpaddr, Handler: handler}
    {{- end }}

    {{- if .HasGRPC }}
    // 6b. Create gRPC server with interceptors
    var grpcServerOpts []grpc.ServerOption
    grpcServerOpts = append(grpcServerOpts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
    grpcServerOpts = append(grpcServerOpts, grpc.ChainUnaryInterceptor(
        log.UnaryServerInterceptor(ctx),
        debug.UnaryServerInterceptor(),
    ))
    grpcServerOpts = append(grpcServerOpts, grpc.ChainStreamInterceptor(
        log.StreamServerInterceptor(ctx),
        debug.StreamServerInterceptor(),
    ))
    grpcServer := grpc.NewServer(grpcServerOpts...)

    {{- range .Services }}
    {{- if .HasGRPC }}
    // {{ .Name }} gRPC server
    {{ .SvcVar }}GRPCServer := {{ .GenGRPCPkg }}.New({{ .EpVar }}, nil)
    {{ .GenGRPCPbPkg }}.Register{{ .StructName }}Server(grpcServer, {{ .SvcVar }}GRPCServer)
    {{- end }}
    {{- end }}

    reflection.Register(grpcServer)
    {{- end }}

    // 7. Start servers (graceful shutdown)
    errc := make(chan error)
    go func() {
        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errc <- fmt.Errorf("%s", <-c)
    }()
    ctx, cancel := context.WithCancel(ctx)

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()

        {{- if .HasHTTP }}
        go func() {
            log.Printf(ctx, "HTTP server listening on %s", *httpaddr)
            errc <- httpServer.ListenAndServe()
        }()
        {{- end }}

        {{- if .HasGRPC }}
        go func() {
            lis, err := net.Listen("tcp", *grpcaddr)
            if err != nil {
                errc <- err
                return
            }
            log.Printf(ctx, "gRPC server listening on %s", *grpcaddr)
            errc <- grpcServer.Serve(lis)
        }()
        {{- end }}

        go func() {
            log.Printf(ctx, "Metrics server listening on %s", *metricsAddr)
            errc <- metricsServer.ListenAndServe()
        }()

        <-ctx.Done()
        log.Printf(ctx, "shutting down servers")

        // Shutdown gracefully with a 30s timeout.
        sctx, scancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer scancel()

        {{- range .Services }}
        if stopper, ok := interface{}({{ .SvcVar }}).(interface{ Stop(context.Context) error }); ok {
            if err := stopper.Stop(sctx); err != nil {
                log.Errorf(sctx, err, "failed to stop service")
            }
        }
        {{- end }}

        {{- if .HasHTTP }}
        if err := httpServer.Shutdown(sctx); err != nil {
            log.Errorf(sctx, err, "failed to shutdown HTTP server")
        }
        {{- end }}

        {{- if .HasGRPC }}
        grpcServer.GracefulStop()
        {{- end }}

        if err := metricsServer.Shutdown(sctx); err != nil {
            log.Errorf(sctx, err, "failed to shutdown metrics server")
        }
    }()

    // Cleanup
    if err := <-errc; err != nil {
        log.Errorf(ctx, err, "exiting")
    }
    cancel()
    wg.Wait()
    log.Printf(ctx, "exited")
}
