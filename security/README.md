# Security Plugin for goa

The `security` package defines [goa
v2](https://github.com/goadesign/goa/tree/v2) DSL functions that make it
possible to describe how the designed services perform auth. The plugin also
contains code generators that take advantage of that description to produce
server and client code as well as the corresponding OpenAPI specification.

The generated server code decodes the security credentials from the incoming
requests and makes them available to user code. The code also initializes the
request context with the security requirements described in the design in the
form of scopes so that authorization code may validate the decoded credentials
against them.

## Security DSL

The security requirements are described in the DSL by first defining the set of
security schemes that are used by the services and then by specifying which
endpoints make use of them.
