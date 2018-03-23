// Code generated with goa v2.0.0-wip, DO NOT EDIT.
//
// secured_service HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/plugins/security/examples/multi_auth/design

package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	securedservice "goa.design/plugins/security/examples/multi_auth/gen/secured_service"
)

// BuildSigninPayload builds the payload for the secured_service signin
// endpoint from CLI flags.
func BuildSigninPayload(securedServiceSigninBody string) (*securedservice.SigninPayload, error) {
	var err error
	var body SigninRequestBody
	{
		err = json.Unmarshal([]byte(securedServiceSigninBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"password\": \"password\",\n      \"username\": \"user\"\n   }'")
		}
	}
	if err != nil {
		return nil, err
	}
	v := &securedservice.SigninPayload{
		Username: body.Username,
		Password: body.Password,
	}
	return v, nil
}

// BuildSecurePayload builds the payload for the secured_service secure
// endpoint from CLI flags.
func BuildSecurePayload(securedServiceSecureBody string, securedServiceSecureFail string) (*securedservice.SecurePayload, error) {
	var err error
	var body SecureRequestBody
	{
		err = json.Unmarshal([]byte(securedServiceSecureBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"token\": \"Dignissimos reiciendis itaque enim quibusdam.\"\n   }'")
		}
	}
	var fail *bool
	{
		if securedServiceSecureFail != "" {
			val, err := strconv.ParseBool(securedServiceSecureFail)
			fail = &val
			if err != nil {
				err = fmt.Errorf("invalid value for fail, must be BOOL")
			}
		}
	}
	if err != nil {
		return nil, err
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
	var err error
	var body DoublySecureRequestBody
	{
		err = json.Unmarshal([]byte(securedServiceDoublySecureBody), &body)
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
	if err != nil {
		return nil, err
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
	var err error
	var body AlsoDoublySecureRequestBody
	{
		err = json.Unmarshal([]byte(securedServiceAlsoDoublySecureBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"oauth_token\": \"Possimus ab asperiores quae deleniti molestiae et.\",\n      \"password\": \"password\",\n      \"token\": \"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ\",\n      \"username\": \"user\"\n   }'")
		}
	}
	var key *string
	{
		if securedServiceAlsoDoublySecureKey != "" {
			key = &securedServiceAlsoDoublySecureKey
		}
	}
	if err != nil {
		return nil, err
	}
	v := &securedservice.AlsoDoublySecurePayload{
		Username:   body.Username,
		Password:   body.Password,
		Token:      body.Token,
		OauthToken: body.OauthToken,
	}
	v.Key = key
	return v, nil
}
