package types

import (
	"fmt"
	"path/filepath"
	"strings"

	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/goa/v3/expr"
)

// GenerateProto produces protobuf message definitions for all user types.
func GenerateProto(genpkg string, roots []eval.Root, files []*codegen.File) ([]*codegen.File, error) {
	for _, root := range roots {
		if r, ok := root.(*expr.RootExpr); ok {
			if f := protoFile(genpkg, r); f != nil {
				files = append(files, f)
			}
		}
	}
	return files, nil
}

func protoFile(genpkg string, r *expr.RootExpr) *codegen.File {
	if len(r.Types) == 0 {
		return nil
	}

	path := filepath.Join(codegen.Gendir, "types", "types.proto")

	var buf strings.Builder
	buf.WriteString("syntax = \"proto3\";\n\n")
	buf.WriteString("package types;\n\n")
	buf.WriteString("option go_package = \"")
	buf.WriteString(genpkg)
	buf.WriteString("/types\";\n\n")

	// Check if we need any imports
	needsAny := false

	for _, t := range r.Types {
		if hasAnyType(t.Attribute()) {
			needsAny = true
		}
		// Future: check for timestamp types
	}

	if needsAny {
		buf.WriteString("import \"google/protobuf/any.proto\";\n\n")
	}

	// Generate messages for each type
	seen := make(map[string]struct{})
	for _, t := range r.Types {
		generateProtoMessage(&buf, t, seen, 0)
	}

	return &codegen.File{
		Path: path,
		SectionTemplates: []*codegen.SectionTemplate{
			{
				Name:   "proto-content",
				Source: "{{ . }}",
				Data:   buf.String(),
			},
		},
	}
}

func generateProtoMessage(buf *strings.Builder, ut expr.UserType, seen map[string]struct{}, indent int) {
	if _, ok := seen[ut.ID()]; ok {
		return
	}
	seen[ut.ID()] = struct{}{}

	// Handle result types and aliases
	if expr.IsAlias(ut) {
		return // Aliases don't generate separate proto messages
	}

	att := ut.Attribute()
	obj := expr.AsObject(att.Type)
	if obj == nil {
		return
	}

	indentStr := strings.Repeat("  ", indent)

	// Write comment if description exists
	if att.Description != "" {
		for _, line := range strings.Split(att.Description, "\n") {
			buf.WriteString(indentStr)
			buf.WriteString("// ")
			buf.WriteString(line)
			buf.WriteString("\n")
		}
	}

	// Write message definition
	buf.WriteString(indentStr)
	buf.WriteString("message ")
	buf.WriteString(ut.Name())
	buf.WriteString(" {\n")

	// Generate fields
	fieldNum := 1
	for _, nat := range *obj {
		generateProtoField(buf, nat, &fieldNum, indent+1)
	}

	buf.WriteString(indentStr)
	buf.WriteString("}\n\n")

	// Recursively generate nested types
	for _, nat := range *obj {
		collectNestedTypes(nat.Attribute, seen, func(nested expr.UserType) {
			generateProtoMessage(buf, nested, seen, indent)
		})
	}
}

func generateProtoField(buf *strings.Builder, nat *expr.NamedAttributeExpr, fieldNum *int, indent int) {
	indentStr := strings.Repeat("  ", indent)

	// Write field comment
	if nat.Attribute.Description != "" {
		buf.WriteString(indentStr)
		buf.WriteString("// ")
		buf.WriteString(nat.Attribute.Description)
		buf.WriteString("\n")
	}

	buf.WriteString(indentStr)

	// Determine optional/repeated modifier
	isRequired := false
	if nat.Attribute.Validation != nil {
		if req := nat.Attribute.Validation.Required; len(req) > 0 {
			for _, r := range req {
				if r == nat.Name {
					isRequired = true
					break
				}
			}
		}
	}

	// Handle repeated (arrays)
	if arr := expr.AsArray(nat.Attribute.Type); arr != nil {
		buf.WriteString("repeated ")
		buf.WriteString(protoType(arr.ElemType.Type))
	} else if m := expr.AsMap(nat.Attribute.Type); m != nil {
		// Maps in proto3
		buf.WriteString("map<")
		buf.WriteString(protoType(m.KeyType.Type))
		buf.WriteString(", ")
		buf.WriteString(protoType(m.ElemType.Type))
		buf.WriteString(">")
	} else {
		// Optional fields in proto3
		if !isRequired && !expr.IsPrimitive(nat.Attribute.Type) {
			buf.WriteString("optional ")
		}
		buf.WriteString(protoType(nat.Attribute.Type))
	}

	buf.WriteString(" ")
	buf.WriteString(codegen.Goify(nat.Name, false))
	buf.WriteString(" = ")
	buf.WriteString(fmt.Sprintf("%d", *fieldNum))
	buf.WriteString(protoJSONOption(nat.Attribute))
	buf.WriteString(";\n")

	*fieldNum++
}

func protoJSONOption(att *expr.AttributeExpr) string {
	if att == nil || att.Meta == nil {
		return ""
	}
	if names := att.Meta["proto:tag:json"]; len(names) > 0 && names[0] != "" {
		return fmt.Sprintf(" [json_name = %q]", names[0])
	}
	return ""
}

func protoType(dt expr.DataType) string {
	switch dt.Kind() {
	case expr.BooleanKind:
		return "bool"
	case expr.IntKind, expr.Int32Kind:
		return "int32"
	case expr.Int64Kind:
		return "int64"
	case expr.UIntKind, expr.UInt32Kind:
		return "uint32"
	case expr.UInt64Kind:
		return "uint64"
	case expr.Float32Kind:
		return "float"
	case expr.Float64Kind:
		return "double"
	case expr.StringKind:
		return "string"
	case expr.BytesKind:
		return "bytes"
	case expr.AnyKind:
		return "google.protobuf.Any"
	case expr.ArrayKind:
		arr := expr.AsArray(dt)
		return protoType(arr.ElemType.Type)
	case expr.MapKind:
		// Maps are handled specially in generateProtoField
		return "map"
	case expr.UserTypeKind:
		ut := dt.(expr.UserType)
		if expr.IsAlias(ut) {
			return protoType(ut.Attribute().Type)
		}
		return ut.Name()
	case expr.ObjectKind:
		// Anonymous objects - could generate inline messages
		return "string" // Fallback
	default:
		return "string" // Safe fallback
	}
}

func collectNestedTypes(att *expr.AttributeExpr, seen map[string]struct{}, cb func(expr.UserType)) {
	if att == nil {
		return
	}

	switch dt := att.Type.(type) {
	case *expr.Object:
		for _, nat := range *dt {
			collectNestedTypes(nat.Attribute, seen, cb)
		}
	case *expr.Array:
		collectNestedTypes(dt.ElemType, seen, cb)
	case *expr.Map:
		collectNestedTypes(dt.KeyType, seen, cb)
		collectNestedTypes(dt.ElemType, seen, cb)
	case expr.UserType:
		if _, ok := seen[dt.ID()]; !ok {
			cb(dt)
		}
	}
}

func hasAnyType(att *expr.AttributeExpr) bool {
	if att == nil {
		return false
	}

	if att.Type.Kind() == expr.AnyKind {
		return true
	}

	switch dt := att.Type.(type) {
	case *expr.Object:
		for _, nat := range *dt {
			if hasAnyType(nat.Attribute) {
				return true
			}
		}
	case *expr.Array:
		return hasAnyType(dt.ElemType)
	case *expr.Map:
		return hasAnyType(dt.KeyType) || hasAnyType(dt.ElemType)
	case expr.UserType:
		return hasAnyType(dt.Attribute())
	}

	return false
}
