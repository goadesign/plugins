# goakit

The `goakit` plugin generates the service code using the [go-kit](https://github.com/go-kit/kit) constructs. In particular:

* The `gen` command generates go-kit decoders and encoders for both the servers and clients.
* The `example` command generates a server that leverages the go-kit HTTP `Server` struct.

## Usage

Simply import the plugin in the service design package:

```go
import _ "goa.design/plugins/goakit"
```

## Example

The [cellar](https://github.com/goa.design/plugins/goakit/examples/cellar)
example
[design](https://github.com/goa.design/plugins/goakit/examples/cellar/design/design.go)
illustrates the usage by importing both the official [cellar example design
package](https://github.com/goa.design/goa/tree/v2/examples/cellar/design) and
the goakit plugin. The go-kit transport code is generated under the
`gen/http/<service name>/kitserver|kitclient` directories.

## Further Use Cases

The generated endpoint functions can be wrapped individually using the various go-kit instrumentation functions, for example:

```go
	// Adapted from the cellar example
	var (
		sommelierEndpoints *sommelier.Endpoints
	)
	{
		sommelierEndpoints = sommelier.NewEndpoints(sommelierSvc)
		sommelierEndpoints.Pick = opentracing.TraceServer(trace, "pick")(sommelierEndpoints.Pick)
		sommelierEndpoints.Pick = LoggingMiddleware(log.With(logger, "method", "Pick"))(sommelierEndpoints.Pick)
	}

```
