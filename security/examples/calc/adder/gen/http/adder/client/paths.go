// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// HTTP request path constructors for the adder service.
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package client

import (
	"fmt"
)

// AddAdderPath returns the URL path to the adder service add HTTP endpoint.
func AddAdderPath(a int, b int) string {
	return fmt.Sprintf("/add/%v/%v", a, b)
}
