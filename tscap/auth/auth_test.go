package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchesAction(t *testing.T) {
	cases := []struct {
		name     string
		granted  []string
		required string
		want     bool
	}{
		{"wildcard grants all", []string{"*"}, "read", true},
		{"exact match", []string{"read"}, "read", true},
		{"no match", []string{"write"}, "read", false},
		{"multiple grants with match", []string{"write", "read"}, "read", true},
		{"multiple grants no match", []string{"write", "delete"}, "read", false},
		{"empty grants", []string{}, "read", false},
		{"wildcard in list", []string{"write", "*"}, "read", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchesAction(tc.granted, tc.required)
			if got != tc.want {
				t.Errorf("matchesAction(%v, %q) = %v, want %v",
					tc.granted, tc.required, got, tc.want)
			}
		})
	}
}

func TestMatchesResource(t *testing.T) {
	cases := []struct {
		name     string
		granted  []string
		required string
		want     bool
	}{
		{"wildcard grants all", []string{"*"}, "items/123", true},
		{"exact match", []string{"items/123"}, "items/123", true},
		{"no match", []string{"items/456"}, "items/123", false},
		{"pattern exact match", []string{"items/*"}, "items/*", true},
		{"multiple grants with match", []string{"users/*", "items/*"}, "items/*", true},
		{"empty grants", []string{}, "items/123", false},
		{"wildcard in list", []string{"users/*", "*"}, "items/123", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := matchesResource(tc.granted, tc.required)
			if got != tc.want {
				t.Errorf("matchesResource(%v, %q) = %v, want %v",
					tc.granted, tc.required, got, tc.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	cases := []struct {
		name       string
		caps       Capabilities
		req        Requirement
		wantOK     bool
		wantStatus int
	}{
		{
			name: "matching capability",
			caps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"read"}, Resources: []string{"*"}},
				},
			},
			req:    Requirement{Capability: "example.com/cap/app", Action: "read", Resource: "items/123"},
			wantOK: true,
		},
		{
			name: "wildcard action",
			caps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"*"}, Resources: []string{"items/*"}},
				},
			},
			req:    Requirement{Capability: "example.com/cap/app", Action: "write", Resource: "items/*"},
			wantOK: true,
		},
		{
			name: "missing capability",
			caps: Capabilities{
				"example.com/cap/other": []Grant{
					{Action: []string{"*"}, Resources: []string{"*"}},
				},
			},
			req:        Requirement{Capability: "example.com/cap/app", Action: "read", Resource: "*"},
			wantOK:     false,
			wantStatus: http.StatusForbidden,
		},
		{
			name: "wrong action",
			caps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"read"}, Resources: []string{"*"}},
				},
			},
			req:        Requirement{Capability: "example.com/cap/app", Action: "write", Resource: "*"},
			wantOK:     false,
			wantStatus: http.StatusForbidden,
		},
		{
			name: "wrong resource",
			caps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"*"}, Resources: []string{"users/*"}},
				},
			},
			req:        Requirement{Capability: "example.com/cap/app", Action: "read", Resource: "items/*"},
			wantOK:     false,
			wantStatus: http.StatusForbidden,
		},
		{
			name: "multiple grants first matches",
			caps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"read"}, Resources: []string{"items/*"}},
					{Action: []string{"write"}, Resources: []string{"users/*"}},
				},
			},
			req:    Requirement{Capability: "example.com/cap/app", Action: "read", Resource: "items/*"},
			wantOK: true,
		},
		{
			name: "multiple grants second matches",
			caps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"read"}, Resources: []string{"items/*"}},
					{Action: []string{"write"}, Resources: []string{"users/*"}},
				},
			},
			req:    Requirement{Capability: "example.com/cap/app", Action: "write", Resource: "users/*"},
			wantOK: true,
		},
		{
			name:       "empty capabilities",
			caps:       Capabilities{},
			req:        Requirement{Capability: "example.com/cap/app", Action: "read", Resource: "*"},
			wantOK:     false,
			wantStatus: http.StatusForbidden,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			got := Check(w, tc.caps, tc.req)
			if got != tc.wantOK {
				t.Errorf("Check() = %v, want %v", got, tc.wantOK)
			}
			if !tc.wantOK && w.Code != tc.wantStatus {
				t.Errorf("Check() status = %d, want %d", w.Code, tc.wantStatus)
			}
		})
	}
}

func TestParseCapabilities(t *testing.T) {
	cases := []struct {
		name       string
		header     string
		wantOK     bool
		wantStatus int
		wantCaps   Capabilities
	}{
		{
			name:   "valid capabilities",
			header: `{"example.com/cap/app":[{"action":["read"],"resources":["*"]}]}`,
			wantOK: true,
			wantCaps: Capabilities{
				"example.com/cap/app": []Grant{
					{Action: []string{"read"}, Resources: []string{"*"}},
				},
			},
		},
		{
			name:       "missing header",
			header:     "",
			wantOK:     false,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid json",
			header:     "not json",
			wantOK:     false,
			wantStatus: http.StatusUnauthorized,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set(Header, tc.header)
			}
			w := httptest.NewRecorder()

			caps, ok := ParseCapabilities(w, req)
			if ok != tc.wantOK {
				t.Errorf("ParseCapabilities() ok = %v, want %v", ok, tc.wantOK)
			}
			if !tc.wantOK && w.Code != tc.wantStatus {
				t.Errorf("ParseCapabilities() status = %d, want %d", w.Code, tc.wantStatus)
			}
			if tc.wantOK {
				if len(caps) != len(tc.wantCaps) {
					t.Errorf("ParseCapabilities() caps length = %d, want %d", len(caps), len(tc.wantCaps))
				}
			}
		})
	}
}
