# goakit Cellar Example

This example demonstrates the use of the goakit plugin to generate a
[go-kit](https://github.com/go-kit/kit) service and client from the goa
[cellar](https://github.com/goadesign/goa/tree/v2/examples/cellar) example.

The only non generated code in the directory is the `design` package which
simply imports the cellar example design and the goakit plugin. Importing the
plugin package is all that is required to generate the go-kit service and
client using the usual `goa` tool commands:

```bash
goa gen goa.design/plugins/goakit/examples/cellar/design
```

The command above generates the content of the `gen` directory. This directory
contains the files that are re-generated each time the `goa` tool is run
reflecting the latest changes in the design.

```bash
goa example goa.design/plugins/goakit/examples/cellar/design
```

The command above generates the content of the `cmd` directory and the top level
example service files (`sommelier.go` and `storage.go`). These files are
generated as an example and won't be overridden by the `goa` tool on future
invocations, instead the tool produces unique filenames each time.

## Code Structure

The generated code contains:

* The service definitions under `gen/sommelier` and `gen/storage`. These define
  the business logic layer in the form of interfaces that gets implemented by
  user code.

* The endpoint definitions under the same directories define endpoints for each
  service methods. Endpoints are transport agnostic call sites that define the
  methods used by clients to call the service and by the service to implement
  the server.

* The HTTP transport under `gen/http/sommelier` and `gen/http/storage`. There
  are a `server` and a `client` packages as well as a `kitserver` and
  `kitclient` packages that wrap the encoder and decoder functions in the
  corresponding non-kit packages into go-kit compatible ones.

* The OpenAPI specification for both services in `gen/http/openapi.json`

* A set of command line parsing functions that can return a endpoint and HTTP
  payloads given command line flags in `gen/http/cli`.

## Composability

The data structures that define the endpoints (e.g.
[sommelier.Endpoints](https://github.com/goadesign/plugins/blob/master/goakit/examples/cellar/gen/sommelier/endpoints.go#L17-L21))
is public and each endpoint method is individually accessible as a struct field
so that user code may wrap it in middleware or completely override it with
custom code.

The same is true of the HTTP server structs (e.g.
[Server](https://github.com/goadesign/plugins/blob/master/goakit/examples/cellar/gen/http/sommelier/server/server.go#L19-L22))
which make it possible to wrap or override each handler individually. The HTTP
client structs (e.g.
[Client](https://github.com/goadesign/plugins/blob/master/goakit/examples/cellar/gen/http/sommelier/client/client.go#L18-L31))
also expose each endpoint client individually (`PickDoer` in this example).
