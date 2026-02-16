# OpenTelemetry Plugin (Deprecated)

> **Deprecated**: This plugin is no longer necessary and will be removed in a
> future release. The Goa HTTP muxer now sets `r.Pattern` on every matched
> request (using the Go 1.22+ convention), which `otelhttp` v0.65.0+ reads
> automatically to tag spans and metrics with the matched route.
>
> **Migration**: Remove the blank import from your design package and
> regenerate:
>
> ```diff
> - import _ "goa.design/plugins/v3/otel"
> ```
>
> No other changes are needed — route tagging now happens automatically in
> the Goa muxer.

## Background

The `otel` plugin was a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that wrapped generated HTTP handlers with `otelhttp.WithRouteTag` to set the
`http.route` attribute on OpenTelemetry spans and metrics.

`otelhttp.WithRouteTag` was
[removed](https://github.com/open-telemetry/opentelemetry-go-contrib/pull/8268)
in `otelhttp` v0.65.0 because `otelhttp` now reads `r.Pattern` (added in Go
1.22) to obtain the route automatically. Goa's muxer has been updated to set
`r.Pattern` on every dispatched request, making this plugin unnecessary.
