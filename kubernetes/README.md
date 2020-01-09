# Kubernetes Plugin

The `kubernetes` plugin is a [Goa v3](https://github.com/goadesign/goa/tree/v3)
plugin that makes it possible to generate deployable kubernetes configuration
yamls for the Goa services.

## Enabling the Plugin

To enable the plugin and make use of the Kubernetes DSLs simply import both the
`kubernetes` and the `dsl` packages as follows

```go
import (
  k8s "goa.design/plugins/v3/kubernetes/dsl"
  . "goa.design/goa/v3/dsl"
)
```

## Effects on Code Generation

Enabling the plugin results in generation of kubernetes configuration yaml
files in the top-level gen directory under `k8s` folder. A yaml file is
is generated for each Goa service that defines the `kubernetes` DSL.

## Design

The `kubernetes` plugin adds a number of functions to the Goa DSL.

Here is an example that shows the usage of kubernetes DSL to generate
kubernetes `Service` and `Deployment` configurations for a Goa service.

```go
var _ = Service("myService", func() {
  k8s.Deployment("my-k8s-deployment", func() {
    k8s.Pod("my-k8s-pod", func() {
      k8s.Replicas(4)
      k8s.Container("bootstrap", k8s.InitContainerKind, func() {
        k8s.Image("org/calc:latest", func() {
          k8s.Command("bootstrap")
        })
      })
      k8s.Container("my-k8s-container", func() {
        k8s.Image("org/calc:latest", func() {
          k8s.Command("foo", "exec")
          k8s.Args("-i", "bar", "-v")
        })
        k8s.Envs(func() {
          k8s.Env("HOST", "foo.bar.svc")
          k8s.Env("PORT", 80)
        })
      })
    })
  })

  k8s.Service("my-k8s-service", func() {
    k8s.Ports(func() {
      k8s.Port(80, "application", 8080, k8s.HTTPKind)
      k8s.Port(9090, "metrics")
    })
  })

  Method("myMethod1", func() {
    ...
  })
  Method("myMethod2", func() {
    ...
  })
})
```

The above design generates a new folder named `k8s` in the top-level gen
directory and contains a file named `myService.yaml` which contains working
kubernetes `Deployment` and `Service` configurations. The yaml file can be
modified as needed and deployed to kubernetes.
