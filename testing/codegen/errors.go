package codegen

import (
	"fmt"
	"maps"
	"path/filepath"
	"slices"
	"text/template"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/service"
	"goa.design/goa/v3/expr"
)

type (
	// errorsData contains data for generating error test helpers.
	errorsData struct {
		Service *service.Data
		Errors  []*errorData
	}

	// errorData describes a service error.
	errorData struct {
		Name        string
		Ref         string
		Description string
		Methods     []string
	}
)

// generateErrors generates error test helpers for a service.
func generateErrors(genpkg string, svcData *service.Data, root *expr.RootExpr, svc *expr.ServiceExpr) *codegen.File {
	if !hasErrors(svc) {
		return nil
	}

	path := filepath.Join(testingPath(genpkg, svc), "errors.go")

	if svcData == nil {
		return nil
	}

	data := buildErrorsData(svcData, svc)

	specs := []*codegen.ImportSpec{
		{Path: "strings"},
		{Path: "testing"},
	}

	sections := []*codegen.SectionTemplate{
		codegen.Header(fmt.Sprintf("Error test helpers for %s service", svc.Name), codegen.SnakeCase(svc.Name)+"test", specs),
		{
			Name:   "error-helpers",
			Source: testingTemplates.Read(errorHelpersT),
			FuncMap: template.FuncMap{
				"goify": codegen.Goify,
			},
			Data: data,
		},
	}

	return &codegen.File{
		Path:             path,
		SectionTemplates: sections,
	}
}

// buildErrorsData builds the template data for error helpers.
func buildErrorsData(svcData *service.Data, svc *expr.ServiceExpr) *errorsData {
	data := &errorsData{
		Service: svcData,
		Errors:  make([]*errorData, 0),
	}

	// Collect all unique errors from methods
	errorMap := make(map[string]*errorData)

	for _, m := range svc.Methods {
		for _, e := range m.Errors {
			key := e.Name
			if existing, found := errorMap[key]; found {
				existing.Methods = append(existing.Methods, m.Name)
			} else {
				ed := &errorData{
					Name:        e.Name,
					Ref:         "", // no concrete type assertion required
					Description: e.Description,
					Methods:     []string{m.Name},
				}
				errorMap[key] = ed
			}
		}
	}

	// Convert map to slice
	for _, k := range slices.Sorted(maps.Keys(errorMap)) {
		data.Errors = append(data.Errors, errorMap[k])
	}

	return data
}

// hasErrors checks if the service has any errors defined.
func hasErrors(svc *expr.ServiceExpr) bool {
	for _, m := range svc.Methods {
		if len(m.Errors) > 0 {
			return true
		}
	}
	return false
}
