package auth

import (
	"encoding/json"
	"net/http"
)

const Header = "Tailscale-App-Capabilities"

// Capabilities represents parsed app capabilities from the Tailscale header.
// The map key is the capability name (e.g., "example.com/cap/myapp").
type Capabilities map[string][]Grant

// Grant represents a single permission grant from a Tailscale ACL.
type Grant struct {
	Action    []string `json:"action"`
	Resources []string `json:"resources"`
}

// Requirement specifies what capability is needed for a method.
type Requirement struct {
	Capability string
	Action     string
	Resource   string
}

// Gate stores the authorization requirements for a method.
type Gate struct {
	MethodName    string
	AllowAnonymous bool
	Requirement   *Requirement
}

// ParseCapabilities extracts capabilities from the Tailscale header.
func ParseCapabilities(w http.ResponseWriter, r *http.Request) (Capabilities, bool) {
	header := r.Header.Get(Header)
	if header == "" {
		WriteUnauthenticated(w, "missing capabilities header")
		return nil, false
	}

	var caps Capabilities
	if err := json.Unmarshal([]byte(header), &caps); err != nil {
		WriteUnauthenticated(w, "invalid capabilities header")
		return nil, false
	}

	return caps, true
}

// Check verifies the caller has the required capability.
func Check(w http.ResponseWriter, caps Capabilities, req Requirement) bool {
	grants, ok := caps[req.Capability]
	if !ok {
		WriteUnauthorized(w, "missing required capability")
		return false
	}

	for _, g := range grants {
		if matchesAction(g.Action, req.Action) && matchesResource(g.Resources, req.Resource) {
			return true
		}
	}

	WriteUnauthorized(w, "insufficient permissions")
	return false
}

// matchesAction checks if the granted actions satisfy the required action.
func matchesAction(granted []string, required string) bool {
	for _, a := range granted {
		if a == "*" || a == required {
			return true
		}
	}
	return false
}

// matchesResource checks if the granted resources satisfy the required resource.
func matchesResource(granted []string, required string) bool {
	for _, r := range granted {
		if r == "*" || r == required {
			return true
		}
	}
	return false
}

// WriteUnauthenticated writes a 401 response.
func WriteUnauthenticated(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unauthenticated",
		"message": message,
	})
}

// WriteUnauthorized writes a 403 response.
func WriteUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   "unauthorized",
		"message": message,
	})
}
