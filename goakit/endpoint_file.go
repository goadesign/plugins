package goakit

import (
	"regexp"

	"goa.design/goa/codegen"
	"goa.design/goa/codegen/service"
	"goa.design/goa/design"
)

// EndpointFile returns the endpoint file for the given service.
func EndpointFile(svc *design.ServiceExpr) *codegen.File {
	f := service.EndpointFile(svc)
	goakitify(f)
	return f
}

// goaEndpointRegexp matches occurrences of the "goa.Endpoint" type in Go code.
var goaEndpointRegexp = regexp.MustCompile(`([^\p{L}_])goa\.Endpoint([^\p{L}_])`)

// goakitify replaces all occurrences of goa.Endpoint with endpoint.Endpoint in
// the file section template sources. It also adds
// "github.com/go-kit/kit/endpoint" to the list of imported packages if
// occurrences were replaces. goakitify modifies the given files and returns
// them.
func goakitify(fs ...*codegen.File) []*codegen.File {
	for _, f := range fs {
		var hasEndpoint bool
		for _, s := range f.SectionTemplates {
			if !hasEndpoint {
				hasEndpoint = goaEndpointRegexp.MatchString(s.Source)
			}
			s.Source = goaEndpointRegexp.ReplaceAllString(s.Source, "${1}endpoint.Endpoint${2}")
		}
		if hasEndpoint {
			codegen.AddImport(
				f.SectionTemplates[0],
				&codegen.ImportSpec{Path: "github.com/go-kit/kit/endpoint"},
			)
		}
	}
	return fs
}
