package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
	modelcodegen "goa.design/model/codegen"
	"goa.design/model/mdl"
	modelexpr "goa.design/plugins/v3/model/expr"
)

// init registers the plugin generator function.
func init() {
	codegen.RegisterPlugin("model", "gen", nil, Validate)
}

// Validate computes the model from the model package path specified in the Root
// expression.  It then validates that each service has a corresponding
// container in the model (and vice versa if ModelComplete was called).
func Validate(_ string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	var services []*expr.ServiceExpr
	var r *modelexpr.RootExpr
	var ok bool
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			// Note: this root is guaranteed to be appear before the model root
			services = r.Services
			continue
		}
		r, ok = root.(*modelexpr.RootExpr)
		if ok {
			break
		}
	}
	if !ok || r.ModelPkgPath == "" {
		return nil, fmt.Errorf("model plugin requires a model package path, use the 'Model' function to set it")
	}
	if r.SystemName == "" {
		return nil, fmt.Errorf("model plugin requires a system name, use the 'Model' function to set it")
	}

	format := r.ContainerNameFormat
	if format == "" {
		format = modelexpr.DefaultFormat
	}
	excluded := r.ExcludedTags
	for i, tag := range excluded {
		excluded[i] = strings.ToLower(tag)
	}

	// Retrieve model JSON
	b, err := modelcodegen.JSON(r.ModelPkgPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to compile model %s: %w", r.ModelPkgPath, err)
	}
	var design mdl.Design
	if err := json.Unmarshal(b, &design); err != nil {
		return nil, fmt.Errorf("failed to load design: %w", err)
	}
	if design.Model == nil {
		return nil, fmt.Errorf("model not found in design")
	}
	var system *mdl.SoftwareSystem
	for _, s := range design.Model.Systems {
		if s.Name == r.SystemName {
			system = s
			break
		}
	}
	if system == nil {
		return nil, fmt.Errorf("system %q not found", r.SystemName)
	}

	// Collect services that do not have a container defined in the model
	var noContainer []string
	for _, svc := range services {
		found := false
		for _, container := range system.Containers {
			var name string
			if n, ok := r.ServiceContainer[svc.Name]; ok {
				if n == "" {
					found = true // ModelNone
					break
				}
				name = n
			} else {
				name = fmt.Sprintf(format, svc.Name)
			}
			if container.Name == name {
				found = true
				break
			}
		}
		if !found {
			noContainer = append(noContainer, svc.Name)
		}
	}

	var noService []string
	if r.Complete {
		// Collect containers that have no corresponding service
	outer:
		for _, container := range system.Containers {
			for _, t := range strings.Split(container.Tags, ",") {
				t = strings.TrimSpace(strings.ToLower(t))
				for _, tag := range excluded {
					if t == tag {
						continue outer
					}
				}
			}
			found := false
			for _, svc := range services {
				for svcName, cont := range r.ServiceContainer {
					if svcName != svc.Name {
						continue
					}
					if cont == container.Name {
						found = true
						break
					}
				}
				if container.Name == fmt.Sprintf(format, svc.Name) {
					found = true
					break
				}
			}
			if !found {
				noService = append(noService, container.Name)
			}
		}
	}

	if len(noContainer) > 0 && len(noService) > 0 {
		return nil, fmt.Errorf(
			"service%s %s %s no corresponding container in the model, and container%s %s %s no corresponding service in the design",
			pluralize(len(noContainer)), strings.Join(noContainer, ", "), pluralVerb(len(noContainer)),
			pluralize(len(noService)), strings.Join(noService, ", "), pluralVerb(len(noService)))
	}
	if len(noContainer) > 0 {
		return nil, fmt.Errorf(
			"service%s %s %s no corresponding container in the model",
			pluralize(len(noContainer)), strings.Join(noContainer, ", "), pluralVerb(len(noContainer)))
	}
	if len(noService) > 0 {
		return nil, fmt.Errorf(
			"container%s %s %s no corresponding service in the design",
			pluralize(len(noService)), strings.Join(noService, ", "), pluralVerb(len(noService)))
	}

	return nil, nil
}

// Helper functions to handle singular/plural
func pluralize(count int) string {
	if count == 1 {
		return ""
	}
	return "s"
}

func pluralVerb(count int) string {
	if count == 1 {
		return "has"
	}
	return "have"
}
