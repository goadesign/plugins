package codegen

import (
	"fmt"

	"goa.design/goa/codegen"
	goadesign "goa.design/goa/design"
	"goa.design/plugins/security/design"
)

type (
	// SchemeData describes a single security scheme.
	SchemeData struct {
		// UsernameField is the name of the payload field that should be
		// initialized with the basic auth username if any.
		UsernameField string
		// UsernamePointer is true if the username field is a pointer.
		UsernamePointer bool
		// PasswordField is the name of the payload field that should be
		// initialized with the basic auth password if any.
		PasswordField string
		// PasswordPointer is true if the password field is a pointer.
		PasswordPointer bool
		// CredField contains the name of the payload field that should
		// be initialized with the API key, the JWT token or the OAuth2
		// access token.
		CredField string
		// CredPointer is true if the credential field is a pointer.
		CredPointer bool
		// KeyAttr is the attribute name containing the security tag.
		KeyAttr string
		// Scheme is the security scheme expression.
		Scheme *design.SchemeExpr
	}
)

// BuildSchemeData builds the scheme data for the given scheme and method expressions.
func BuildSchemeData(s *design.SchemeExpr, m *goadesign.MethodExpr) *SchemeData {
	if !goadesign.IsObject(m.Payload.Type) {
		return nil
	}
	switch s.Kind {
	case design.BasicAuthKind:
		userAtt, user := findSecurityField(m.Payload, "security:username")
		passAtt, pass := findSecurityField(m.Payload, "security:password")
		return &SchemeData{
			UsernameField:   user,
			UsernamePointer: m.Payload.IsPrimitivePointer(userAtt, true),
			PasswordField:   pass,
			PasswordPointer: m.Payload.IsPrimitivePointer(passAtt, true),
			Scheme:          s,
		}
	case design.APIKeyKind:
		if keyAtt, key := findSecurityField(m.Payload, "security:apikey:"+s.SchemeName); key != "" {
			return &SchemeData{
				CredField:   key,
				CredPointer: m.Payload.IsPrimitivePointer(keyAtt, true),
				KeyAttr:     keyAtt,
				Scheme:      s,
			}
		}
	case design.JWTKind:
		if keyAtt, key := findSecurityField(m.Payload, "security:token"); key != "" {
			return &SchemeData{
				CredField:   key,
				CredPointer: m.Payload.IsPrimitivePointer(keyAtt, true),
				KeyAttr:     keyAtt,
				Scheme:      s,
			}
		}
	case design.OAuth2Kind:
		if keyAtt, key := findSecurityField(m.Payload, "security:accesstoken"); key != "" {
			return &SchemeData{
				CredField:   key,
				CredPointer: m.Payload.IsPrimitivePointer(keyAtt, true),
				KeyAttr:     keyAtt,
				Scheme:      s,
			}
		}
	default:
		panic(fmt.Sprintf("unknown kind %#v", s.Kind)) // bug
	}
	return nil
}

// findSecurityField returns the name and corresponding field name of child
// attribute of p with the given tag if p is an object.
func findSecurityField(a *goadesign.AttributeExpr, tag string) (string, string) {
	obj := goadesign.AsObject(a.Type)
	if obj == nil {
		return "", ""
	}
	for _, at := range *obj {
		if _, ok := at.Attribute.Metadata[tag]; ok {
			return at.Name, codegen.Goify(at.Name, true)
		}
	}
	return "", ""
}
