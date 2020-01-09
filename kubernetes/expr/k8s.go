package expr

import (
	"fmt"

	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

type (
	// ServiceExpr describes the kubernetes service.
	ServiceExpr struct {
		// Name is the service name.
		Name string
		// Ports is the list of ports exposed by the service.
		Ports []*PortExpr
		// Service is the Goa service expression.
		Service *expr.ServiceExpr
	}

	// DeploymentExpr describes the kubernetes deployment.
	DeploymentExpr struct {
		// Name is the deployment name.
		Name string
		// Pod is the kubernetes pod managed by the deployment.
		Pod *PodExpr
		// Service is the Goa service expression.
		Service *expr.ServiceExpr
	}

	// PodExpr describes the kubernetes pod.
	PodExpr struct {
		// Name is the pod name.
		Name string
		// Replicas is the number of pod replicas (defaults to 1).
		Replicas int
		// Containers is the list of containers run by the pod.
		Containers []*ContainerExpr
	}

	// ContainerExpr describes the kubernetes container.
	ContainerExpr struct {
		// Name is the container name.
		Name string
		// Kind is the kind of container (regular or init).
		Kind ContainerKind
		// Image is the image used by the container.
		Image *ImageExpr
		// Ports is the list of ports exposed by the container.
		Ports []*PortExpr
		// Envs is the list of environment variables set in the container.
		Envs []*EnvExpr
	}

	// ImageExpr describes the image used in the kubernetes container.
	ImageExpr struct {
		// Name is the image name.
		Name string
		// Command is the command to be run in the image.
		Command []string
		// Args is the arguments to the command.
		Args []string
	}

	// PortExpr describes the port exposed by the kubernetes service or
	// container.
	PortExpr struct {
		// Value is the port number.
		Value int
		// Name is the port name.
		Name string
		// Protocol is the protocol used by the port.
		Protocol ProtocolKind
		// Target is the target port for a kubernetes service or container port
		// for a kubernetes container.
		Target int
		// Parent is the parent expression of the port (service or container).
		Parent eval.Expression
	}

	// EnvExpr describes the environment variable set in a kubernetes container.
	EnvExpr struct {
		// Name is the variable name.
		Name string
		// Value is the variable value.
		Value interface{}
	}

	// ContainerKind is a type to indicate the type of the kubernetes container.
	ContainerKind int

	// ProtocolKind is a type to indicate the type of protocol used by a
	// kubernetes port.
	ProtocolKind string
)

const (
	// RegularContainerKind indicates that the container is a regular container.
	RegularContainerKind ContainerKind = iota + 1
	// InitContainerKind indicates that the container is an init container.
	InitContainerKind
)

const (
	// TCPKind indicates the protocol type is TCP.
	TCPKind ProtocolKind = "TCP"
	// UDPKind indicates the protocol type is UDP.
	UDPKind ProtocolKind = "UDP"
	// HTTPKind indicates the protocol type is HTTP.
	HTTPKind ProtocolKind = "HTTP"
	// ProxyProtocolKind indicates the protocol type is Proxy protocol.
	ProxyProtocolKind ProtocolKind = "Proxy"
	// SCTPKind indicates the protocol type is SCTP.
	SCTPKind ProtocolKind = "SCTP"
)

// EvalName is the name of service expression.
func (s *ServiceExpr) EvalName() string {
	return fmt.Sprintf("kubernetes service %q", s.Name)
}

// Prepare prepares the service expression.
func (s *ServiceExpr) Prepare() {
	if len(s.Ports) == 0 {
		s.Ports = append(s.Ports, &PortExpr{Value: 80})
	}
}

// Finalize finalizes service ports.
func (s *ServiceExpr) Finalize() {
	for _, p := range s.Ports {
		p.finalize()
	}
}

// EvalName is the name of deployment expression.
func (d *DeploymentExpr) EvalName() string {
	return fmt.Sprintf("kubernetes deployment %q", d.Name)
}

// Validate validates the deployment expression.
func (d *DeploymentExpr) Validate() error {
	verr := new(eval.ValidationErrors)
	if d.Pod == nil {
		verr.Add(d, "No pod set in deployment")
	} else {
		verr.Merge(d.Pod.validate())
	}
	return verr
}

// Finalize finalizes deployment pod.
func (d *DeploymentExpr) Finalize() {
	d.Pod.finalize()
}

// EvalName is the name of container expression.
func (c *ContainerExpr) EvalName() string {
	return fmt.Sprintf("kubernetes container %q", c.Name)
}

func (c *ContainerExpr) validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	if c.Image == nil {
		verr.Add(c, "No image set for the container")
	}
	return verr
}

// finalize finalizes pod containers.
func (c *ContainerExpr) finalize() {
	for _, p := range c.Ports {
		p.finalize()
	}
}

// EvalName is the name of pod expression.
func (p *PodExpr) EvalName() string {
	return fmt.Sprintf("kubernetes pod %q", p.Name)
}

// validate validates the pod expression.
func (p *PodExpr) validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	clen := len(p.Containers)
	switch {
	case clen == 0:
		verr.Add(p, "No container defined in pod")
	case clen > 1:
		// validate container names are unique in a pod
		for i := 0; i < len(p.Containers)-1; i++ {
			c1 := p.Containers[i]
			for _, c2 := range p.Containers[i+1:] {
				if c1.Name == c2.Name {
					verr.Add(p, "Multiple containers in pod have same name %q", c1.Name)
					break
				}
			}
		}
	}
	// validate containers
	for _, c := range p.Containers {
		verr.Merge(c.validate())
	}
	return verr
}

// finalize finalizes pod containers.
func (p *PodExpr) finalize() {
	for _, c := range p.Containers {
		c.finalize()
	}
}

// EvalName is the name of image expression.
func (i *ImageExpr) EvalName() string {
	return fmt.Sprintf("kubernetes container image %q", i.Name)
}

// EvalName is the name of port expression.
func (p *PortExpr) EvalName() string {
	n := fmt.Sprintf("%d_%s", p.Value, p.Protocol)
	if p.Name != "" {
		n = fmt.Sprintf("%s_%s", n, p.Name)
	}
	return fmt.Sprintf("kubernetes port %q", n)
}

func (p *PortExpr) validate() *eval.ValidationErrors {
	verr := new(eval.ValidationErrors)
	if p.Value < 1 || p.Value > 65535 {
		verr.Add(p, "Invalid port number. Port number must be between 0 and 65536.")
	}
	switch t := p.Parent.(type) {
	case *ContainerExpr:
		if p.Protocol != TCPKind && p.Protocol != UDPKind {
			verr.Add(p, "Invalid port type %q for container %q. Port type must be TCP or UDP for container.", p.Protocol, t.Name)
		}
	}
	return verr
}

func (p *PortExpr) finalize() {
	if p.Protocol == "" {
		p.Protocol = TCPKind
	}
	if p.Target == 0 {
		p.Target = p.Value
	}
}
