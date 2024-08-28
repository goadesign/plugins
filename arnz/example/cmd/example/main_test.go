package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

const (
	admin           = "arn:aws:iam::123456789012:user/administrator"
	dev             = "arn:aws:iam::123456789012:user/developer"
	readonly        = "arn:aws:iam::123456789012:user/read-only"
	created         = "{\"action\":\"created!\"}\n"
	read            = "{\"action\":\"read!\"}\n"
	updated         = "{\"action\":\"updated!\"}\n"
	deleted         = "{\"action\":\"deleted!\"}\n"
	healthy         = "{\"action\":\"healthy!\"}\n"
	unauthenticated = "{\"error\":\"unauthenticated\",\"message\":\"caller not authenticated\"}\n"
	unauthorized    = "{\"error\":\"unauthorized\",\"message\":\"caller not authorized\"}\n"
)

func TestUnsigned(t *testing.T) {
	cases := []struct {
		Name       string
		Verb       string
		Path       string
		StatusCode int
		Body       string
	}{
		{"Unsigned Create", "POST", "/", http.StatusUnauthorized, unauthenticated},
		{"Unsigned Read", "GET", "/", http.StatusUnauthorized, unauthenticated},
		{"Unsigned Update", "PUT", "/", http.StatusUnauthorized, unauthenticated},
		{"Unsigned Delete", "DELETE", "/", http.StatusUnauthorized, unauthenticated},
		{"Unsigned Health", "GET", "/health", http.StatusOK, healthy},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ts := httptest.NewServer(server(8080).Handler)
			defer ts.Close()

			resp := unsigned(t, c.Verb, ts.URL+c.Path)
			assert.Equal(t, c.StatusCode, resp.StatusCode)

			bytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, c.Body, string(bytes))
		})
	}
}

func TestSigned(t *testing.T) {
	cases := []struct {
		Name       string
		Caller     string
		Verb       string
		Path       string
		StatusCode int
		Body       string
	}{
		{"Admin Create", admin, "POST", "/", http.StatusOK, created},
		{"Admin Read", admin, "GET", "/", http.StatusOK, read},
		{"Admin Update", admin, "PUT", "/", http.StatusOK, updated},
		{"Admin Delete", admin, "DELETE", "/", http.StatusOK, deleted},
		{"Admin Health", admin, "GET", "/health", http.StatusOK, healthy},

		{"Dev Create", dev, "POST", "/", http.StatusForbidden, unauthorized},
		{"Dev Read", dev, "GET", "/", http.StatusOK, read},
		{"Dev Update", dev, "PUT", "/", http.StatusOK, updated},
		{"Dev Delete", dev, "DELETE", "/", http.StatusForbidden, unauthorized},
		{"Dev Health", dev, "GET", "/health", http.StatusOK, healthy},

		{"ReadOnly Create", readonly, "POST", "/", http.StatusForbidden, unauthorized},
		{"ReadOnly Read", readonly, "GET", "/", http.StatusOK, read},
		{"ReadOnly Update", readonly, "PUT", "/", http.StatusForbidden, unauthorized},
		{"ReadOnly Delete", readonly, "DELETE", "/", http.StatusForbidden, unauthorized},
		{"ReadOnly Health", readonly, "GET", "/health", http.StatusOK, healthy},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ts := httptest.NewServer(server(8080).Handler)
			defer ts.Close()

			resp := signed(t, c.Verb, ts.URL+c.Path, c.Caller)
			assert.Equal(t, c.StatusCode, resp.StatusCode)

			bytes, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, c.Body, string(bytes))
		})
	}
}

func unsigned(t *testing.T, verb, url string) *http.Response {
	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		t.Fatalf("Failed to create %s request: %v", verb, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make %s request: %v", verb, err)
	}

	return resp
}

func signed(t *testing.T, verb, url, callerArn string) *http.Response {
	amznReqCtx := events.APIGatewayV2HTTPRequestContext{
		Authorizer: &events.APIGatewayV2HTTPRequestContextAuthorizerDescription{
			IAM: &events.APIGatewayV2HTTPRequestContextAuthorizerIAMDescription{
				UserARN: callerArn,
			},
		},
	}

	header, err := json.Marshal(amznReqCtx)
	if err != nil {
		t.Fatalf("Failed to marshal header: %v", err)
	}

	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Amzn-Request-Context", string(header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	return resp
}
