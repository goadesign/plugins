// Package mains provides an example-phase plugin that rewrites the example
// server mains to follow the Goa "golden path" conventions. It relocates
// single-service servers under services/<svc>/cmd/<svc>/main.go, and keeps
// multi-service servers under cmd/<server>/main.go, wiring Clue/OTel, debug
// endpoints, and health/metrics consistently.
package mains

import (
    "path"
    "path/filepath"
    "strings"

    "goa.design/goa/v3/codegen"
    "goa.design/goa/v3/codegen/example"
    "goa.design/goa/v3/codegen/service"
    "goa.design/goa/v3/eval"
    "goa.design/goa/v3/expr"
    httpcodegen "goa.design/goa/v3/http/codegen"
)

const (
    // pluginName is the registered plugin name.
    pluginName = "mains"
    // pluginCmd is the goa CLI command the plugin integrates with.
    pluginCmd  = "example"
)

// srvInfo stores server-level data derived from generated example files.
type srvInfo struct {
    Dir        string
    APIPkg     string
    Services   []*service.Data
    HasWS      bool
    ServerName string
}

// svcT provides template data for each service imported by a server.
type svcT struct {
    Name         string
    StructName   string
    SvcVar       string
    EpVar        string
    SrvVar       string
    GenPkg       string
    GenHTTPPkg   string
    HasWebSocket bool
}

// Register the plugin for the example phase.
func init() {
    codegen.RegisterPluginLast(pluginName, pluginCmd, nil, Generate)
}

// Generate produces golden-path mains that follow the Goa conventions and
// pulse weather layout:
//  - For servers with a single service: services/<svc>/cmd/<svc>/main.go
//  - For servers with multiple services: cmd/<server>/main.go
// It replaces the default example main and http.go files.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
    return generateExample(genpkg, roots, files)
}

func generateExample(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
    // Collect per-server services/APIPkg from example main files first
    srvMap := map[string]*srvInfo{}
    for _, f := range files {
        if filepath.Base(f.Path) != "main.go" { continue }
        segs := strings.Split(filepath.ToSlash(f.Path), "/")
        if len(segs) < 3 || segs[0] != "cmd" { continue }
        dir := segs[1]
        var svcs []*service.Data
        var apipkg string
        for _, s := range f.SectionTemplates {
            switch s.Name {
            case "server-main-services":
                if dm, ok := s.Data.(map[string]any); ok {
                    if v, ok := dm["Services"].([]*service.Data); ok { svcs = v }
                }
            case "server-main-logger":
                if dm, ok := s.Data.(map[string]any); ok {
                    if v, ok := dm["APIPkg"].(string); ok { apipkg = v }
                }
            }
        }
        if len(svcs) == 0 { continue }
        if apipkg == "" { apipkg = apiPkgAlias(genpkg, roots) }
        if _, exists := srvMap[dir]; !exists {
            srvMap[dir] = &srvInfo{Dir: dir, APIPkg: apipkg, Services: svcs}
        }
    }
    // Complement with HTTP server data (WebSocket detection) and fallback if needed
    for _, f := range files {
        if filepath.Base(f.Path) != "http.go" { continue }
        segs := strings.Split(filepath.ToSlash(f.Path), "/")
        if len(segs) < 3 || segs[0] != "cmd" { continue }
        dir := segs[1]
        var httpSvcs []*httpcodegen.ServiceData
        for _, s := range f.SectionTemplates {
            if s.Name == "server-http-start" {
                if dm, ok := s.Data.(map[string]any); ok {
                    if v, ok := dm["Services"].([]*httpcodegen.ServiceData); ok { httpSvcs = v; break }
                }
            }
        }
        if len(httpSvcs) == 0 { continue }
        var svcs []*service.Data
        for _, sd := range httpSvcs { if sd != nil && sd.Service != nil { svcs = append(svcs, sd.Service) } }
        hasWS := httpcodegen.NeedDialer(httpSvcs)
        apipkg := apiPkgAlias(genpkg, roots)
        if info, ok := srvMap[dir]; ok {
            info.HasWS = hasWS
            if info.APIPkg == "" { info.APIPkg = apipkg }
        } else {
            srvMap[dir] = &srvInfo{Dir: dir, APIPkg: apipkg, Services: svcs, HasWS: hasWS}
        }
    }

    if len(srvMap) == 0 {
        return files, nil
    }

    // Filter out default example mains and http.go; we'll add our own mains.
    var out []*codegen.File
    for _, f := range files {
        base := filepath.Base(f.Path)
        if strings.HasPrefix(f.Path, "cmd/") && (base == "main.go" || base == "http.go") {
            continue
        }
        out = append(out, f)
    }

    // Create mains per server
    for _, info := range srvMap {
        if len(info.Services) == 0 {
            continue
        }
        specs := []*codegen.ImportSpec{
            {Path: "context"},
            {Path: "flag"},
            {Path: "fmt"},
            {Path: "net/http"},
            {Path: "net/http/httptrace"},
            {Path: "os"},
            {Path: "os/signal"},
            {Path: "sync"},
            {Path: "syscall"},
            {Path: "time"},
            {Path: "go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"},
            {Path: "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"},
            {Path: "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"},
            {Path: "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"},
            {Path: "goa.design/clue/clue"},
            {Path: "goa.design/clue/debug"},
            {Path: "goa.design/clue/health"},
            {Path: "goa.design/clue/log"},
            codegen.GoaNamedImport("http", "goahttp"),
            {Path: "google.golang.org/grpc/credentials/insecure"},
        }
        if info.HasWS {
            specs = append(specs, &codegen.ImportSpec{Path: "github.com/gorilla/websocket"})
        }
        rootPath := moduleRootFromGenpkg(genpkg)
        specs = append(specs, &codegen.ImportSpec{Path: rootPath, Name: info.APIPkg})

        scope := codegen.NewNameScope()
        var svcsData []svcT
        wsBySvc := httpWebSocketByService(roots)
        hasAnyWS := false
        for _, sd := range info.Services {
            genAlias := scope.Unique(sd.PkgName, "svc")
            httpAlias := scope.Unique(sd.PkgName+"svr", "svr")
            specs = append(specs,
                &codegen.ImportSpec{Path: path.Join(genpkg, sd.PathName), Name: genAlias},
                &codegen.ImportSpec{Path: path.Join(genpkg, "http", sd.PathName, "server"), Name: httpAlias},
            )
            hws := wsBySvc[sd.Name]
            if hws {
                hasAnyWS = true
            }
            svcsData = append(svcsData, svcT{
                Name:         sd.Name,
                StructName:   sd.StructName,
                SvcVar:       sd.VarName + "Svc",
                EpVar:        sd.VarName + "Endpoints",
                SrvVar:       sd.VarName + "Server",
                GenPkg:       genAlias,
                GenHTTPPkg:   httpAlias,
                HasWebSocket: hws,
            })
        }

        sections := []*codegen.SectionTemplate{
            codegen.Header("", "main", specs),
            {Name: "mains-main", Source: tmpl.Read("main"), Data: map[string]any{
                "APIPkg":          info.APIPkg,
                "Services":        svcsData,
                "HasAnyWebSocket": hasAnyWS,
                "ServiceCount":    len(svcsData),
                "ServerLabel":     serverLabel(roots),
            }},
        }

        var fpath string
        if len(info.Services) == 1 {
            svc := info.Services[0]
            fpath = filepath.ToSlash(filepath.Join("services", svc.PathName, "cmd", svc.PathName, "main.go"))
        } else {
            // Use the example server directory so people can still `go run ./cmd/<dir>` when multiple services
            // are served from a single process.
            // Derive server directory using same logic as example generator.
            fpath = filepath.ToSlash(filepath.Join("cmd", example.Servers.Get(rootServer(roots), roots[0].(*expr.RootExpr)).Dir, "main.go"))
            // But we already filtered out the default http.go/main.go, so only our main remains.
        }

        out = append(out, &codegen.File{Path: fpath, SectionTemplates: sections, SkipExist: true})
    }
    return out, nil
}

func apiPkgAlias(genpkg string, roots []eval.Root) string {
    var apiName string
    for _, r := range roots {
        if root, ok := r.(*expr.RootExpr); ok {
            if root != nil && root.API != nil {
                apiName = root.API.Name
                break
            }
        }
    }
    if apiName == "" {
        apiName = "api"
    }
    scope := codegen.NewNameScope()
    return scope.Unique(strings.ToLower(codegen.Goify(apiName, false)), "api")
}

func serverLabel(roots []eval.Root) string {
    for _, r := range roots {
        if root, ok := r.(*expr.RootExpr); ok {
            if root != nil && root.API != nil {
                return strings.ToLower(codegen.Goify(root.API.Name, false))
            }
        }
    }
    return "goa-service"
}

func moduleRootFromGenpkg(genpkg string) string {
    idx := strings.LastIndex(genpkg, "/")
    if idx <= 0 {
        return "."
    }
    return genpkg[:idx]
}

func httpWebSocketByService(roots []eval.Root) map[string]bool {
    hasWS := map[string]bool{}
    for _, r := range roots {
        root, ok := r.(*expr.RootExpr)
        if !ok || root.API == nil || root.API.HTTP == nil {
            continue
        }
        for _, svc := range root.API.HTTP.Services {
            for _, e := range svc.HTTPEndpoints {
                if e.SSE != nil {
                    continue
                }
                if e.MethodExpr != nil && e.MethodExpr.Stream != expr.NoStreamKind {
                    hasWS[svc.Name()] = true
                    break
                }
            }
        }
    }
    return hasWS
}

// rootServer returns the first server expression if any.
func rootServer(roots []eval.Root) *expr.ServerExpr {
    for _, r := range roots {
        if root, ok := r.(*expr.RootExpr); ok {
            if root.API != nil && len(root.API.Servers) > 0 {
                return root.API.Servers[0]
            }
        }
    }
    return nil
}
