package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	cases := []struct {
		Name       string
		Path       string
		StatusCode int
	}{
		{"AuthO", "/autho", http.StatusForbidden},
		{"AuthZ", "/authz", http.StatusForbidden},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ts := httptest.NewServer(server(8080).Handler)
			defer ts.Close()

			resp, err := http.Get(ts.URL + c.Path)
			if err != nil {
				t.Fatalf("Failed to make GET request: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, c.StatusCode, resp.StatusCode)
		})
	}
}
