// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// adder HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/examples/calc/adder/design

package client

import (
	"fmt"
	"strconv"

	addersvc "goa.design/plugins/security/examples/calc/adder/gen/adder"
)

// BuildAddAddPayload builds the payload for the adder add endpoint from CLI
// flags.
func BuildAddAddPayload(adderAddA string, adderAddB string, adderAddKey string) (*addersvc.AddPayload, error) {
	var err error
	var a int
	{
		var v int64
		v, err = strconv.ParseInt(adderAddA, 10, 64)
		a = int(v)
		if err != nil {
			err = fmt.Errorf("invalid value for a, must be INT")
		}
	}
	var b int
	{
		var v int64
		v, err = strconv.ParseInt(adderAddB, 10, 64)
		b = int(v)
		if err != nil {
			err = fmt.Errorf("invalid value for b, must be INT")
		}
	}
	var key string
	{
		key = adderAddKey
	}
	if err != nil {
		return nil, err
	}
	payload := &addersvc.AddPayload{
		A:   a,
		B:   b,
		Key: key,
	}
	return payload, nil
}
