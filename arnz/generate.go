package arnz

import (
	"path/filepath"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	goahttp "goa.design/goa/v3/http/codegen"
)

var MethodARNs = make(map[string][]string)

func init() {
	codegen.RegisterPlugin("arnz", "gen", nil, Generate)
}

type ArnzData struct {
	MethodName  string
	AllowedArns []string
}

func Generate(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, file := range files {
		if filepath.Base(file.Path) == "server.go" {
			for _, s := range file.Section("server-handler") {
				if data, ok := s.Data.(*goahttp.EndpointData); ok {
					if _, ok := MethodARNs[data.Method.Name]; ok {
						codegen.AddImport(file.SectionTemplates[0],
							&codegen.ImportSpec{Path: "encoding/json"},
							&codegen.ImportSpec{Path: "strings"},
							&codegen.ImportSpec{Path: "github.com/aws/aws-lambda-go/events"})

						file.SectionTemplates = append(file.SectionTemplates, &codegen.SectionTemplate{
							Name:   "arnz-middleware",
							Source: arnzMiddleWareT,
							Data: ArnzData{
								MethodName:  data.Method.Name,
								AllowedArns: MethodARNs[data.Method.Name],
							},
						})

						s.Source = serverHandlerT
					}
				}
			}
		}
	}
	return files, nil
}

const arnzMiddleWareT = `
{{ printf "for authorization based on AWS ARNs" | comment }}
func {{ .MethodName }}Arnz(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var authorized bool
		key := "X-Amzn-Request-Context"

		authorized = false
		allowedArns := []string{ 
			{{- range .AllowedArns }} 
			{{ printf "%q" . }}, 
			{{- end }}
		}

		amzReqCtxHeader := r.Header.Get(key)

		var amzCtx events.APIGatewayV2HTTPRequestContext{}
		err := json.Unmarshal([]byte(amzReqCtxHeader), &amzCtx)
		if err != nil {
			http.Error(w, "unauthorized: error parsing X-Amzn-Request-Context header", http.StatusUnauthorized)
			return
		}

		if amzCtx.Authorizer == nil {
			http.Error(w, "unauthorized: no authorizer defined in X-Amzn-Request-Context", http.StatusUnauthorized)
			return
		}

		if amzCtx.Authorizer.IAM == nil {
			http.Error(w, "unauthorized: no IAM authorizer defined in X-Amzn-Request-Context", http.StatusUnauthorized)
			return
		}

		callerArn := amzCtx.Authorizer.IAM.UserARN
		for _, allowedArn := range allowedArns {
			if strings.Contains(callerArn, allowedArn) {
				authorized = true
				break
			}
		}

		if !authorized {
			http.Error(w, 
				fmt.Sprintf("unauthorized: %s", callerArn), 
				http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
`

const serverHandlerT = `
{{ printf "%s configures the mux to serve the %q service %q endpoint." .MountHandler .ServiceName .Method.Name | comment }}
func {{ .MountHandler }}(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	{{- $methodName := .Method.Name }}
	{{- range .Routes }}
	mux.Handle("{{ .Verb }}", "{{ .Path }}", {{ $methodName }}Arnz(f))
	{{- end }}
}
`
