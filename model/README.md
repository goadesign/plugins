# Model Plugin

The `model` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3) plugin
that validates a [Model](https://github.com/goadesign/model) diagram DSL against
a Goa design. The plugin generates errors if there are services present in the
Goa design with no corresponding container in the Model diagram. The plugin can
also validate that all the containers in the Model diagram have a corresponding
service in the Goa design.

## Enabling the Plugin

To enable the plugin simply import the model plugin `dsl` package in the Goa
design:

```go
import (
  . "goa.design/goa/v3/dsl"
  . "goa.design/plugins/v3/model/dsl"
)
```

Running `goa gen` will now execute the plugin and validate that all services
have a corresponding container in the model (unless configured otherwise, see
below).

## Usage

The plugin extends the Goa DSL as follows:

* The `Model` function sets the `model` Go package and the name of the software
  system to validate against, this is the only required DSL.
* The `ModelContainerNameFormat` function sets the format used by the plugin to
  map service to container names.
* The `ModelExcludedTags` function sets a list of case-insensitive tags that
  exclude containers equipped with them from the validation process.
* The `ModelComplete` function enables the validation that all the containers in
  the Model diagram have a corresponding service in the Goa design.

All these functions must appear inside the `API` function:

```go
var _ = API("calc", func() {
  // Set the model package and the name of the software system
  Model("goa.design/model/examples/basic/model", "calc")

  // Set the container name format, default is "%s Service"
  ModelContainerFormat("%s Service")

  // Exclude containers with the "database", "external" and "thirdparty" tags (default)
  ModelExcludedTags("database", "external", "thirdparty")

  // Enable the validation that all the containers in the Model diagram have a
  // corresponding service. Not enabled by default.
  ModelComplete()
})
```

The plugin also defines the following DSL functions that must appear under a
`Service` definition:

* The `ModelContainer` function sets the name of the corresponding diagram
  container. This overrides the name computed from the `ContainerNameFormat`
  format.
* The `ModelNone` function disables the validation for the service.

```go
var _ = Service("add", func() {
  // Set the container name
  ModelContainer("Add Service")
  // ...
})

var _ = Service("no diagram", func() {
  // Disable the validation for this service
  ModelNone()
  // ...
})
```