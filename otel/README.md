# OpenTelemetry Plugin

The `otel` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin that instruments
HTTP endpoints with OpenTelemetry configuring it with the endpoint route pattern.

## Usage

Simply import the plugin in the service design package. Use the blank identifier `_` as explicit
package name:

```go
package design

import . "goa.design/goa/v3/dsl"
import _ "goa.design/plugins/v3/otel" # Enables the otel plugin

var _ = API("...
```

and generate as usual:

```bash
goa gen PACKAGE
```

where `PACKAGE` is the Go import path of the design package.

## Effects of the Plugin

Importing the `otel` package changes the behavior of the `gen` command of the
`goa` tool. The `gen` command output is modified so that the generated HTTP
handlers are wrapped with a call to the `otelhttp.WithRouteTag`
