package kubernetes

import (
	"fmt"
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/eval"
	goaexpr "goa.design/goa/v3/expr"
	"goa.design/plugins/v3/kubernetes/expr"
)

type (
	// serviceData contains the data required to generate kubernetes Service
	// config.
	serviceData struct {
		// Name is the service name.
		Name string
		// Ports is the list of ports exposed by the service.
		Ports []*portData
	}

	// deploymentData contains the data required to generate kubernetes Deployment
	// config.
	deploymentData struct {
		// Name is the deployment name.
		Name string
		// Pod represents the kubernetes Pod managed by the deployment.
		Pod *podData
	}

	// podData contains the data required to generate kubernetes Pod template in
	// kubernetes Deployment.
	podData struct {
		// Name is the pod name.
		Name string
		// Replicas is the number of pod replicas to deploy.
		Replicas int
		// Containers is the list of containers in the pod.
		Containers []*containerData
		// InitContainers is the list of init containers in the pod.
		InitContainers []*containerData
	}

	// containerData contains the data required to generate kubernetes Container
	// template in a kubernetes Pod.
	containerData struct {
		// Name is the unique container name.
		Name string
		// Image is the image used by the container.
		Image *imageData
		// Ports is the list of container ports.
		Ports []*portData
		// Envs is the list of environment variables passed to the container.
		Envs []*envData
		// Init, if true, indicates that the container is an init container.
		Init bool
	}

	// imageData contains the data required to generate image spec for a
	// kubernetes container.
	imageData struct {
		// Name is the image name to be pulled.
		Name string
		// Command is the command to be executed.
		Command []string
		// Args is the arguments to the command.
		Args []string
	}

	// portData contains the data required to generate port spec exposed by a
	// kubernetes Service or Container.
	portData struct {
		// Value is the port value.
		Value int
		// Name is the port name.
		Name string
		// Protocol is the protocol used in the port.
		Protocol string
		// Target is the target port exposed in a service.
		Target int
	}

	// envData contains the data required to generate environment variables spec
	// for the kubernetes Container.
	envData struct {
		// Name is the variable name.
		Name string
		// Value is the variable value.
		Value interface{}
	}
)

// Register the plugin Generator functions.
func init() {
	codegen.RegisterPlugin("kubernetes", "gen", nil, Generate)
}

// Generate produces the kubernetes config files for the Goa services that
// define the kuberbetes DSL.
func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		r, ok := root.(*expr.RootExpr)
		if !ok {
			continue
		}
		for _, d := range r.Deployments {
			if f := generateDeploymentConfigFile(genpkg, d); f != nil {
				files = append(files, f)
			}
		}
		for _, s := range r.Services {
			if f := generateServiceConfigFile(genpkg, s); f != nil {
				files = append(files, f)
			}
		}
	}
	return files, nil
}

func generateServiceConfigFile(genpkg string, s *expr.ServiceExpr) *codegen.File {
	data := &serviceData{Name: s.Name}
	for _, p := range s.Ports {
		data.Ports = append(data.Ports, buildPortData(p))
	}
	sections := []*codegen.SectionTemplate{
		&codegen.SectionTemplate{
			Name:   "k8s-service",
			Source: svcConfigT,
			Data:   data,
		},
	}
	return &codegen.File{Path: filePath(s.Service), SectionTemplates: sections}
}

func generateDeploymentConfigFile(genpkg string, d *expr.DeploymentExpr) *codegen.File {
	data := &deploymentData{
		Name: d.Name,
		Pod: &podData{
			Name:     d.Pod.Name,
			Replicas: d.Pod.Replicas,
		},
	}
	for _, c := range d.Pod.Containers {
		cdata := &containerData{
			Name: c.Name,
			Image: &imageData{
				Name:    c.Image.Name,
				Command: c.Image.Command,
				Args:    c.Image.Args,
			},
		}
		for _, p := range c.Ports {
			cdata.Ports = append(cdata.Ports, buildPortData(p))
		}
		for _, e := range c.Envs {
			cdata.Envs = append(cdata.Envs, &envData{Name: e.Name, Value: e.Value})
		}
		if c.Kind == expr.InitContainerKind {
			cdata.Init = true
			data.Pod.InitContainers = append(data.Pod.InitContainers, cdata)
		} else {
			data.Pod.Containers = append(data.Pod.Containers, cdata)
		}
	}
	sections := []*codegen.SectionTemplate{
		&codegen.SectionTemplate{
			Name:   "k8s-deployment",
			Source: deplConfigT,
			Data:   data,
			FuncMap: map[string]interface{}{
				"printList": func(arr []string) string {
					if len(arr) == 0 {
						return ""
					}
					output := fmt.Sprintf("%q", arr[0])
					for _, s := range arr[1:] {
						output += fmt.Sprintf(", %q", s)
					}
					return output
				},
			},
		},
	}
	return &codegen.File{Path: filePath(d.Service), SectionTemplates: sections}
}

func filePath(s *goaexpr.ServiceExpr) string {
	svc := service.Services.Get(s.Name)
	return filepath.Join(codegen.Gendir, "k8s", fmt.Sprintf("%s.yaml", codegen.SnakeCase(svc.Name)))
}

func buildPortData(p *expr.PortExpr) *portData {
	return &portData{
		Value:    p.Value,
		Name:     p.Name,
		Protocol: string(p.Protocol),
		Target:   p.Target,
	}
}

var (
	// Input: serviceData
	svcConfigT = `---
apiVersion: v1
kind: Service
metadata:
  name: {{ .Name }}
spec:
  ports:
{{- range .Ports }}
  - port: {{ .Value }}
    protocol: {{ .Protocol }}
    targetPort: {{ .Target }}
  {{- if .Name }}
    name: {{ .Name }}
  {{- end }}
{{- end }}
`

	// Input: deploymentData
	deplConfigT = `---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {{ .Name }}
spec:
  replicas: {{ .Pod.Replicas }}
	minReadySeconds: 30
	strategy:
	  rollingUpdate:
		  maxUnavailable: 1
			maxSurge: 0
	selector:
	  matchLabels:
		  app: {{ .Name }}
	template:
	  metadata:
		  labels:
			  app: {{ .Name }}
		spec:
		{{- if gt .Pod.Replicas 1 }}
		  affinity:
			  podAntiAffinity:
				  preferredDuringSchedulingIgnoredDuringExecution:
					  - podAffinityTerm:
						    labelSelector:
								  matchExpressions:
									- key: app
									  operator: In
										values:
										- {{ .Name }}
								topologyKey: failure-domain.beta.kubernetes.io/zone
						  weight: 100
		{{- end }}
	  {{- range .Pod.InitContainers }}
	    initContainers:
	    {{- template "container" . }}
	  {{- end }}
	  {{- range .Pod.Containers }}
	    containers:
	    {{- template "container" . }}
	  {{- end }}

{{- define "container" }}
      - name: {{ .Name }}
		    image: {{ .Image.Name }}
		    imagePullPolicy: IfNotPresent
		  {{- $cmd := (printList .Image.Command) }}
		  {{- if $cmd }}
		    command: [{{ $cmd }}]
			{{- end }}
		  {{- $args := (printList .Image.Args) }}
		  {{- if $args }}
		    args: [{{ $args }}]
			{{- end }}
			{{- if not .Init }}
		    livenessProbe:
		      httpGet:
			      path: /health-check
				    port: 8080
			    initialDelaySeconds: 5
			    periodSeconds: 10
		    readinessProbe:
		      httpGet:
			      path: /health-check
				    port: 8080
			    initialDelaySeconds: 5
			    periodSeconds: 10
			{{- end }}
				resources:
				  requests:
					  memory: "128M"
					limits:
					  memory: "128M"
		{{- if .Ports }}
		    ports:
      {{- range .Ports }}
	      - containerPort: {{ .Value }}
		    {{- if .Name }}
		      name: {{ .Name }}
		    {{- end }}
				  protocol: {{ .Protocol }}
	    {{- end }}
		{{- end }}
	      env:
		    - name: POD_NAME
		      valueFrom:
			      fieldRef:
				      fieldPath: metadata.name
	    {{- range .Envs }}
	      - name: {{ .Name }}
		      value: "{{ .Value }}"
	    {{- end }}
{{- end }}
`
)
