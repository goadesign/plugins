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
func BuildSecurePayload(securedServiceSecureFail string, securedServiceSecureToken string) (*securedservice.SecurePayload, error) {
	var err error
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
	var token *string
	{
		if securedServiceSecureToken != "" {
			token = &securedServiceSecureToken
		}
	}
	if err != nil {
		return nil, err
	}
	payload := &securedservice.SecurePayload{
		Fail:  fail,
		Token: token,
	}
	return payload, nil
}

// BuildDoublySecurePayload builds the payload for the secured_service
// doubly_secure endpoint from CLI flags.
func BuildDoublySecurePayload(securedServiceDoublySecureKey string, securedServiceDoublySecureToken string) (*securedservice.DoublySecurePayload, error) {
	var key *string
	{
		if securedServiceDoublySecureKey != "" {
			key = &securedServiceDoublySecureKey
		}
	}
	var token *string
	{
		if securedServiceDoublySecureToken != "" {
			token = &securedServiceDoublySecureToken
		}
	}
	payload := &securedservice.DoublySecurePayload{
		Key:   key,
		Token: token,
	}
	return payload, nil
}

// BuildAlsoDoublySecurePayload builds the payload for the secured_service
// also_doubly_secure endpoint from CLI flags.
func BuildAlsoDoublySecurePayload(securedServiceAlsoDoublySecureBody string, securedServiceAlsoDoublySecureKey string, securedServiceAlsoDoublySecureToken string, securedServiceAlsoDoublySecureOauthToken string) (*securedservice.AlsoDoublySecurePayload, error) {
	var err error
	var body AlsoDoublySecureRequestBody
	{
		err = json.Unmarshal([]byte(securedServiceAlsoDoublySecureBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, example of valid JSON:\n%s", "'{\n      \"password\": \"password\",\n      \"username\": \"user\"\n   }'")
		}
	}
	var key *string
	{
		if securedServiceAlsoDoublySecureKey != "" {
			key = &securedServiceAlsoDoublySecureKey
		}
	}
	var token *string
	{
		if securedServiceAlsoDoublySecureToken != "" {
			token = &securedServiceAlsoDoublySecureToken
		}
	}
	var oauthToken *string
	{
		if securedServiceAlsoDoublySecureOauthToken != "" {
			oauthToken = &securedServiceAlsoDoublySecureOauthToken
		}
	}
	if err != nil {
		return nil, err
	}
	v := &securedservice.AlsoDoublySecurePayload{
		Username: body.Username,
		Password: body.Password,
	}
	v.Key = key
	v.Token = token
	v.OauthToken = oauthToken
	return v, nil
}
