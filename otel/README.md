# OpenTelemetry Plugin

The `otel` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that sets the `http.route` OpenTelemetry span attribute on generated HTTP
handlers. This ensures the matched route pattern appears as a span attribute
regardless of how `otelhttp` is wired into the application.

## Usage

Import the plugin in the service design package with a blank identifier:

```go
package design

import . "goa.design/goa/v3/dsl"
import _ "goa.design/plugins/v3/otel"

var _ = API("...
```

and generate as usual:

```bash
goa gen PACKAGE
```

The generated `MountXxxHandler` functions will set `http.route` on the active
span before calling the handler:

```go
mux.Handle("GET", "/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    trace.SpanFromContext(r.Context()).SetAttributes(semconv.HTTPRoute("/users/{id}"))
    f(w, r)
})
```

## Alternative: mux middleware (no plugin needed)

As of Goa v3.25.0 the default muxer sets `r.Pattern` before middlewares run.
If you register `otelhttp` as a mux middleware, it reads the route pattern
automatically and this plugin is not necessary:

```go
mux := goahttp.NewMuxer()
mux.Use(otelhttp.NewMiddleware("service"))
```

Use the plugin when you wrap the mux externally with
`otelhttp.NewHandler(mux, ...)`, since in that case `r.Pattern` is not yet
available when `otelhttp` creates the span.
