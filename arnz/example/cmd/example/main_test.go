package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	d "goa.design/plugins/v3/arnz/example/design"
)

func TestUnsigned(t *testing.T) {
	cases := []struct {
		Name       string
		Verb       string
		Path       string
		StatusCode int
	}{
		{"Anonymous Create AllowUnsigned", "POST", "/like", http.StatusUnauthorized},
		{"Anonymous Read AllowUnsigned", "GET", "/like", http.StatusUnauthorized},
		{"Anonymous Update AllowUnsigned", "PUT", "/like", http.StatusOK},
		{"Anonymous Delete AllowUnsigned", "DELETE", "/like", http.StatusUnauthorized},

		{"Anonymous Create AllowUnsigned", "POST", "/match", http.StatusUnauthorized},
		{"Anonymous Read AllowUnsigned", "GET", "/match", http.StatusUnauthorized},
		{"Anonymous Update AllowUnsigned", "PUT", "/match", http.StatusOK},
		{"Anonymous Delete AllowUnsigned", "DELETE", "/match", http.StatusUnauthorized},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ts := httptest.NewServer(server(8080).Handler)
			defer ts.Close()

			resp := unsigned(t, c.Verb, ts.URL+c.Path)
			assert.Equal(t, c.StatusCode, resp.StatusCode)
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
	}{
		{"Signed Create AllowArnsLike AdminArn", d.AdminArn, "POST", "/like", http.StatusOK},
		{"Signed Read AllowArnsLike AdminArn", d.AdminArn, "GET", "/like", http.StatusOK},
		{"Signed Update AllowArnsLike AdminArn", d.AdminArn, "PUT", "/like", http.StatusOK},
		{"Signed Delete AllowArnsLike AdminArn", d.AdminArn, "DELETE", "/like", http.StatusOK},

		{"Signed Create AllowArnsLike DevArn", d.DevArn, "POST", "/like", http.StatusUnauthorized},
		{"Signed Read AllowArnsLike DevArn", d.DevArn, "GET", "/like", http.StatusOK},
		{"Signed Update AllowArnsLike DevArn", d.DevArn, "PUT", "/like", http.StatusOK},
		{"Signed Delete AllowArnsLike DevArn", d.DevArn, "DELETE", "/like", http.StatusUnauthorized},

		{"Signed Create AllowArnsLike ReadArn", d.ReadArn, "POST", "/like", http.StatusUnauthorized},
		{"Signed Read AllowArnsLike ReadArn", d.ReadArn, "GET", "/like", http.StatusOK},
		{"Signed Update AllowArnsLike ReadArn", d.ReadArn, "PUT", "/like", http.StatusUnauthorized},
		{"Signed Delete AllowArnsLike ReadArn", d.ReadArn, "DELETE", "/like", http.StatusUnauthorized},

		{"Signed Create AllowArnsMatching AdminArn", d.AdminArn, "POST", "/match", http.StatusOK},
		{"Signed Read AllowArnsMatching AdminArn", d.AdminArn, "GET", "/match", http.StatusOK},
		{"Signed Update AllowArnsMatching AdminArn", d.AdminArn, "PUT", "/match", http.StatusOK},
		{"Signed Delete AllowArnsMatching AdminArn", d.AdminArn, "DELETE", "/match", http.StatusOK},

		{"Signed Create AllowArnsMatching DevArn", d.DevArn, "POST", "/match", http.StatusUnauthorized},
		{"Signed Read AllowArnsMatching DevArn", d.DevArn, "GET", "/match", http.StatusOK},
		{"Signed Update AllowArnsMatching DevArn", d.DevArn, "PUT", "/match", http.StatusOK},
		{"Signed Delete AllowArnsMatching DevArn", d.DevArn, "DELETE", "/match", http.StatusUnauthorized},

		{"Signed Create AllowArnsMatching ReadArn", d.ReadArn, "POST", "/match", http.StatusUnauthorized},
		{"Signed Read AllowArnsMatching ReadArn", d.ReadArn, "GET", "/match", http.StatusOK},
		{"Signed Update AllowArnsMatching ReadArn", d.ReadArn, "PUT", "/match", http.StatusUnauthorized},
		{"Signed Delete AllowArnsMatching ReadArn", d.ReadArn, "DELETE", "/match", http.StatusUnauthorized},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ts := httptest.NewServer(server(8080).Handler)
			defer ts.Close()

			resp := signed(t, c.Verb, ts.URL+c.Path, c.Caller)
			assert.Equal(t, c.StatusCode, resp.StatusCode)
		})
	}
}

func unsigned(t *testing.T, verb, url string) *http.Response {
	req, err := http.NewRequest(verb, url, nil)
	if err != nil {
		t.Fatalf("Failed to create %s request: %v", verb, err)
	}

	// Set the Content-Type header to application/json
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

	// Set the necessary headers for JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Amzn-Request-Context", string(header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	return resp
}
