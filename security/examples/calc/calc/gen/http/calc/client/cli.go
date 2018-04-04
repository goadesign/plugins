// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// calc HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/calc/design

package client

import (
	"fmt"
	"strconv"

	calcsvc "goa.design/plugins/security/examples/calc/calc/gen/calc"
)

// BuildLoginPayload builds the payload for the calc login endpoint from CLI
// flags.
func BuildLoginPayload(calcLoginUser string, calcLoginPassword string) (*calcsvc.LoginPayload, error) {
	var user string
	{
		user = calcLoginUser
	}
	var password string
	{
		password = calcLoginPassword
	}
	payload := &calcsvc.LoginPayload{
		User:     user,
		Password: password,
	}
	return payload, nil
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
