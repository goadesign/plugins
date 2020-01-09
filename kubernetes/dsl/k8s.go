package dsl

import (
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/kubernetes/expr"

	// Register code generators for the plugin
	_ "goa.design/plugins/v3/kubernetes"
)

const (
	// RegularContainerKind indicates that the kubernetes container is a regular
	// container type.
	RegularContainerKind = expr.RegularContainerKind
	// InitContainerKind indicates that the kubernetes container is an
	// init container.
	InitContainerKind = expr.InitContainerKind

	// TCPKind indicates the protocol type is TCP.
	TCPKind = expr.TCPKind
	// UDPKind indicates the protocol type is UDP.
	UDPKind = expr.UDPKind
	// HTTPKind indicates the protocol type is HTTP.
	HTTPKind = expr.HTTPKind
	// ProxyProtocolKind indicates the protocol type is Proxy protocol.
	ProxyProtocolKind = expr.ProxyProtocolKind
	// SCTPKind indicates the protocol type is SCTP.
	SCTPKind = expr.SCTPKind
)

// Service defines the kubernetes service config.
//
// Service must appear in Service expression.
//
// Service accepts service name string as the first argument and an optional
// DSL function as the second argument.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    var _ = Service("calc", func() {
//        k8s.Service("calc-svc", func() {
//            k8s.Ports(func() {
//                k8s.Port(8080, "http")
//            })
//        })
//    })
//
func Service(name string, fns ...func()) {
	if len(fns) > 1 {
		eval.ReportError("too many arguments given to Service")
	}
	svc, ok := eval.Current().(*goaexpr.ServiceExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	s := &expr.ServiceExpr{Name: name, Service: svc}
	if len(fns) > 0 {
		eval.Execute(fns[0], s)
	}
	expr.Root.Services = append(expr.Root.Services, s)
}

// Deployment defines the kubernetes deployment config.
//
// Deployment must appear in Service expression.
//
// Deployment accepts deployment metadata name as the first argument and a
// required function argument defining the other deployment attributes.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    var _ = Service("calc", func() {
//        k8s.Deployment("calc-depl", func() {
//            k8s.Pod("calc-pod", func() {
//                k8s.Container("calc", func() {
//                    k8s.Image("myorg/calc:latest")
//                })
//            })
//        })
//    })
//
func Deployment(name string, fn func()) {
	svc, ok := eval.Current().(*goaexpr.ServiceExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	d := &expr.DeploymentExpr{Name: name, Service: svc}
	eval.Execute(fn, d)
	expr.Root.Deployments = append(expr.Root.Deployments, d)
}

// Pod defines the kubernetes pod config section.
//
// Pod must appear in kubernetes Deployment expression.
//
// Pod accepts pod name as the first argument and a required function as the
// second argument defining the pod attributes.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Deployment("calc-depl", func() {
//        k8s.Pod("calc-pod", func() {
//            k8s.Container("calc", func() {
//                k8s.Image("myorg/calc:latest")
//            })
//        })
//    })
//
func Pod(name string, fn func()) {
	d, ok := eval.Current().(*expr.DeploymentExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	p := &expr.PodExpr{Name: name, Replicas: 1}
	eval.Execute(fn, p)
	d.Pod = p
}

// Replicas sets the number of kubernetes pod replicas.
//
// Replicas must appear in the kubernetes Pod expression.
//
// Replicas accept one integer indicating the number of pod replicas.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Deployment("calc-depl", func() {
//        k8s.Pod("calc-pod", func() {
//            k8s.Replicas(2)
//            k8s.Container("calc", func() {
//                k8s.Image("myorg/calc:latest")
//            })
//        })
//    })
//
func Replicas(n int) {
	if n <= 0 {
		eval.ReportError("Replicas cannot be an integer less than or equal to zero")
	}
	p, ok := eval.Current().(*expr.PodExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	p.Replicas = n
}

// Container defines the container config section for a kubernetes pod.
//
// Container must appear in kubernetes Pod expression.
//
// Container accepts container name as the first argument, an optional type
// as the second argument, and a function argument defining the container
// attributes. The optional type argument can be either RegularContainerKind or
// InitContainerKind. If optional type is not specified then the container is
// set to RegularContainerKind.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Container("calc", func() {
//        k8s.Image("myorg/calc:latest")
//    })
//
//    k8s.Container("bootstrap", INIT, func() {
//        k8s.Image("myorg/calc:latest")
//        k8s.Envs(func() {
//            k8s.Env("HOST", "myhost")
//            k8s.Env("PORT", 80)
//        })
//    })
//
func Container(name string, args ...interface{}) {
	p, ok := eval.Current().(*expr.PodExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	c := &expr.ContainerExpr{Name: name, Kind: RegularContainerKind}
	var dslfn func()
	switch len(args) {
	case 0:
		eval.ReportError("too few arguments given to Container")
	case 1:
		fn, ok := args[0].(func())
		if !ok {
			eval.InvalidArgError("functionr", args[0])
		}
		dslfn = fn
	case 2:
		k, ok := args[0].(expr.ContainerKind)
		if !ok {
			eval.InvalidArgError("second argument is not a container kind", args[0])
		}
		c.Kind = k
		fn, ok := args[1].(func())
		if !ok {
			eval.InvalidArgError("function", args[1])
		}
		dslfn = fn
	default:
		eval.ReportError("too many arguments given to Container")
	}
	eval.Execute(dslfn, c)
	p.Containers = append(p.Containers, c)
}

// Ports defines the kubernetes ports config section.
//
// Ports must appear in kubernetes Service or Container expression.
//
// Ports accepts one argument which is a function listing the ports.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Service("calc-svc", func() {
//        k8s.Ports(func() {
//            k8s.Port(80, "http")
//            k8s.Port(443, "https")
//        })
//    })
//
func Ports(fn func()) {
	switch actual := eval.Current().(type) {
	case *expr.ServiceExpr, *expr.ContainerExpr:
		eval.Execute(fn, actual)
	default:
		eval.IncompatibleDSL()
	}
}

// Port defines the kubernetes port config section.
//
// Port must appear in Ports expression,
//
// Port accepts an integer port number as the first argument, an optional
// port name as the second argument, an optional protocol kind as the third
// argument, and an optional port number referring to the target port number
// as the fourth argument. The allowed protocol kinds are TCPKind, UDPKind,
// HTTPKind, ProxyKind, and SCTPKind. By default, the protocol is set to
// TCPKind and the target port number is set to the integer port number.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Ports(func() {
//        k8s.Port(80)
//        k8s.Port(8080, "http")
//        k8s.Port(8080, "http", TCPKind)
//        k8s.Port(8443, "https", HTTPKind, 443)
//    })
//
func Port(n int, args ...interface{}) {
	p := &expr.PortExpr{Value: n, Parent: eval.Current()}
	if len(args) > 3 {
		eval.ReportError("too many arguments to Port")
	}
	if len(args) > 0 {
		s, ok := args[0].(string)
		if !ok {
			eval.InvalidArgError("second argument to Port must be a string", args[0])
		}
		p.Name = s
		if len(args) > 1 {
			k, ok := args[1].(expr.ProtocolKind)
			if !ok {
				eval.InvalidArgError("third argument to Port must be a supported ProtocolKind", args[1])
			}
			p.Protocol = k
			if len(args) == 3 {
				i, ok := args[2].(int)
				if !ok {
					eval.InvalidArgError("fourth argument to Port must be an int", args[2])
				}
				p.Target = i
			}
		}
	}
	switch actual := eval.Current().(type) {
	case *expr.ServiceExpr:
		actual.Ports = append(actual.Ports, p)
	case *expr.ContainerExpr:
		actual.Ports = append(actual.Ports, p)
	default:
		eval.IncompatibleDSL()
	}
}

// Image defines the image attributes used by the kubernetes container.
//
// Image must appear in the kubernetes Container expression.
//
// Image acceps image to pull as the first argument and an optional function
// to specify the command and arguments as the second argument.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Container("calc", func() {
//        k8s.Image("myorg/calc:latest", func() {
//            k8s.Command("foo", "bar")
//            k8s.Args("-c", "-v")
//        })
//    })
//
func Image(name string, args ...interface{}) {
	c, ok := eval.Current().(*expr.ContainerExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	if len(args) > 1 {
		eval.ReportError("too many arguments to Image")
	}
	i := &expr.ImageExpr{Name: name}
	if len(args) == 1 {
		fn, ok := args[0].(func())
		if !ok {
			eval.InvalidArgError("function", args[0])
		}
		eval.Execute(fn, i)
	}
	c.Image = i
}

// Command defines the command to run on an image in a kubernetes container.
//
// Command must appear in the kubernetes Image expression.
//
// Command accepts one or more strings as the command to run.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Image("myorg/calc:latest", func() {
//        k8s.Command("foo", "bar")
//    })
//
func Command(cmds ...string) {
	i, ok := eval.Current().(*expr.ImageExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	i.Command = cmds
}

// Args defines the command arguments to pass to command in a kubernetes
// container.
//
// Args must appear in the kubernetes Image expression.
//
// Args accepts one or more strings as the command arguments to pass.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Image("myorg/calc:latest", func() {
//        k8s.Command("foo", "bar")
//        k8s.Args("-c", "-v")
//    })
//
func Args(args ...string) {
	i, ok := eval.Current().(*expr.ImageExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	i.Args = args
}

// Envs defines the environment variables to set in the kubernetes container.
//
// Envs must appear in the kubernetes Container expression.
//
// Envs accepts one argument which is a function listing the environment
// variables.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Container("bootstrap", INIT, func() {
//        k8s.Image("myorg/calc:latest")
//        k8s.Envs(func() {
//            k8s.Env("HOST", "myhost")
//            k8s.Env("PORT", 80)
//        })
//    })
//
func Envs(fn func()) {
	c, ok := eval.Current().(*expr.ContainerExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	eval.Execute(fn, c)
}

// Env defines an environment variable to set in the kubernetes container.
//
// Env must appear in the kubernetes Envs expression.
//
// Env accepts environment variable name as the first argument and its value
// as the second argument.
//
// Example:
//
//    import k8s "goa.design/plugins/v3/kubernetes/dsl"
//
//    k8s.Container("bootstrap", INIT, func() {
//        k8s.Image("myorg/calc:latest")
//        k8s.Envs(func() {
//            k8s.Env("HOST", "myhost")
//            k8s.Env("PORT", 80)
//        })
//    })
//
func Env(name string, args ...interface{}) {
	c, ok := eval.Current().(*expr.ContainerExpr)
	if !ok {
		eval.IncompatibleDSL()
	}
	if len(args) > 1 {
		eval.ReportError("too many arguments to Env")
	}
	c.Envs = append(c.Envs, &expr.EnvExpr{Name: name, Value: args[0]})
}
