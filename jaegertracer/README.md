# JaegerTracer Plugin

The `JaegerTracer` plugin is a [goa v2](https://github.com/goadesign/goa/tree/v2) plugin
that adds initialization code in the generated example main file for a [Jaeger](https://https://github.com/jaegertracing/jaeger-client-go) tracer.

## Enabling the Plugin

The plugin requires that the Tracing expression be used in the API or Server definition.

To enable it, import it in your design.go file using the blank identifier `_` as follows:

```go

package design

import . "goa.design/goa/http/design"
import . "goa.design/goa/http/dsl"
import _ "goa.design/plugins/jaegertracer" # Enables the plugin

var _ = API(

    Tracing("localhost:5775") # Required in API or Server
    ...
)
```

and generate as usual:

```bash
goa gen PACKAGE
goa example PACKAGE
```

where `PACKAGE` is the Go import path of the design package.