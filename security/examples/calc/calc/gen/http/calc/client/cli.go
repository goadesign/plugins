// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// BuildLoginPayload builds the payload for the calc login endpoint from CLI
// flags.
func BuildLoginPayload(calcLoginBody string) (*calcsvc.LoginPayload, error) {
	var err error
	var body LoginRequestBody
	{
		err = json.Unmarshal([]byte(calcLoginBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"password\": \"password\",\n      \"user\": \"username\"\n   }'")
		}
	}
	if err != nil {
		return nil, err
	}
	v := &calcsvc.LoginPayload{
		User:     body.User,
		Password: body.Password,
	}
	return v, nil
}

// BuildAddPayload builds the payload for the calc add endpoint from CLI flags.
func BuildAddPayload(calcAddA string, calcAddB string, calcAddToken string) (*calcsvc.AddPayload, error) {
	var err error
	var a int
	{
		var v int64
		v, err = strconv.ParseInt(calcAddA, 10, 64)
		a = int(v)
		if err != nil {
			err = fmt.Errorf("invalid value for a, must be INT")
		}
	}
	var b int
	{
		var v int64
		v, err = strconv.ParseInt(calcAddB, 10, 64)
		b = int(v)
		if err != nil {
			err = fmt.Errorf("invalid value for b, must be INT")
		}
	}
	var token string
	{
		token = calcAddToken
	}
	if err != nil {
		return nil, err
	}
	payload := &calcsvc.AddPayload{
		A:     a,
		B:     b,
		Token: token,
	}
	return payload, nil
}
