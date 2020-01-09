package design

import (
	. "goa.design/goa/v3/dsl"
	k8s "goa.design/plugins/v3/kubernetes/dsl"
)

var _ = Service("calc", func() {
	k8s.Service("calc", func() {
		k8s.Ports(func() {
			k8s.Port(80, "http")
			k8s.Port(443, "https", k8s.HTTPKind, 8443)
		})
	})
	k8s.Deployment("calc", func() {
		k8s.Pod("calc", func() {
			k8s.Replicas(2)
			k8s.Container("bootstrap", k8s.InitContainerKind, func() {
				k8s.Image("myorg/calc:latest")
			})
			k8s.Container("calc", func() {
				k8s.Image("myorg/calc:latest", func() {
					k8s.Command("foo", "bar", "-c")
					k8s.Args("-k", "-v")
				})
				k8s.Ports(func() {
					k8s.Port(8080)
					k8s.Port(9090, "metrics")
				})
				k8s.Envs(func() {
					k8s.Env("HOSTNAME", "example.com")
					k8s.Env("PORT", 80)
				})
			})
		})
	})
	Method("add", func() {
		Payload(func() {
			Attribute("a", Int, func() {
				Example(1)
			})
			Attribute("b", Int, func() {
				Example(2)
			})
			Required("a", "b")
		})
		Result(Int, func() {
			Example(3)
		})
		HTTP(func() {
			GET("/add/{a}/{b}")

			Response(StatusOK)
		})
	})
})
