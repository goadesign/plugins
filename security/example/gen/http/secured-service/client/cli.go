// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/example/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"goa.design/plugins/security/example/gen/securedservice"
)

// BuildSecurePayload builds the payload for the secured_service secure
// endpoint from CLI flags.
func BuildSecurePayload(securedServiceSecureBody string, securedServiceSecureFail string) (*securedservice.SecurePayload, error) {
	var body SecureRequestBody
	{
		err := json.Unmarshal([]byte(securedServiceSecureBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"token\": \"Voluptatem architecto consequatur fuga nisi veritatis.\"\n   }'")
		}
	}
	var fail *bool
	{
		if securedServiceSecureFail != "" {
			val, err := strconv.ParseBool(securedServiceSecureFail)
			if err != nil {
				return nil, fmt.Errorf("invalid value for fail, must be BOOL")
			}
			fail = &val
		}
	}
	v := &securedservice.SecurePayload{
		Token: body.Token,
	}
	v.Fail = fail

	return v, nil
}

// BuildDoublySecurePayload builds the payload for the secured_service
// doubly_secure endpoint from CLI flags.
func BuildDoublySecurePayload(securedServiceDoublySecureBody string, securedServiceDoublySecureKey string) (*securedservice.DoublySecurePayload, error) {
	var body DoublySecureRequestBody
	{
		err := json.Unmarshal([]byte(securedServiceDoublySecureBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"token\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ\"\n   }'")
		}
	}
	var key *string
	{
		if securedServiceDoublySecureKey != "" {
			key = &securedServiceDoublySecureKey
		}
	}
	v := &securedservice.DoublySecurePayload{
		Token: body.Token,
	}
	v.Key = key

	return v, nil
}

// BuildAlsoDoublySecurePayload builds the payload for the secured_service
// also_doubly_secure endpoint from CLI flags.
func BuildAlsoDoublySecurePayload(securedServiceAlsoDoublySecureBody string, securedServiceAlsoDoublySecureKey string) (*securedservice.AlsoDoublySecurePayload, error) {
	var body AlsoDoublySecureRequestBody
	{
		err := json.Unmarshal([]byte(securedServiceAlsoDoublySecureBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"token\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ\"\n   }'")
		}
	}
	var key *string
	{
		if securedServiceAlsoDoublySecureKey != "" {
			key = &securedServiceAlsoDoublySecureKey
		}
	}
	v := &securedservice.AlsoDoublySecurePayload{
		Token: body.Token,
	}
	v.Key = key

	return v, nil
}
