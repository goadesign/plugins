package docs

import "goa.design/goa/v3/http/codegen/openapi"

type (
	// data is the data structure that is serialized to create the docs.
	data struct {
		API         *apiData                   `json:"api"`
		Services    map[string]*serviceData    `json:"services"`
		Definitions map[string]*openapi.Schema `json:"definitions"`
	}

	apiData struct {
		Name         string                 `json:"name"`
		Title        string                 `json:"title,omitempty"`
		Description  string                 `json:"description,omitempty"`
		Version      string                 `json:"version,omitempty"`
		Servers      map[string]*serverData `json:"servers,omitempty"`
		Terms        string                 `json:"terms,omitempty"`
		Contact      *contactData           `json:"contact,omitempty"`
		License      *licenseData           `json:"license,omitempty"`
		Docs         *docsData              `json:"docs,omitempty"`
		Requirements []*requirementData     `json:"requirements,omitempty"`
	}

	serverData struct {
		Name        string               `json:"name"`
		Description string               `json:"description,omitempty"`
		Services    []string             `json:"services,omitempty"`
		Hosts       map[string]*hostData `json:"hosts,omitempty"`
	}

	contactData struct {
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
		URL   string `json:"url,omitempty"`
	}

	licenseData struct {
		Name string `json:"name"`
		URL  string `json:"url,omitempty"`
	}

	docsData struct {
		Description string `json:"description,omitempty"`
		URL         string `json:"url,omitempty"`
	}

	hostData struct {
		Name        string          `json:"name"`
		ServerName  string          `json:"server,omitempty"`
		Description string          `json:"description,omitempty"`
		URIs        []string        `json:"uris,omitempty"`
		Variables   []*variableData `json:"variables,omitempty"`
	}

	variableData struct {
		Name         string   `json:"name"`
		DefaultValue string   `json:"default,omitempty"`
		Enum         []string `json:"enum,omitempty"`
	}

	serviceData struct {
		Name         string                 `json:"name"`
		Description  string                 `json:"description,omitempty"`
		Methods      map[string]*methodData `json:"methods,omitempty"`
		Requirements []*requirementData     `json:"schemes,omitempty"`
	}

	methodData struct {
		Name         string                `json:"name"`
		Description  string                `json:"description,omitempty"`
		Payload      *payloadData          `json:"payload,omitempty"`
		Result       *payloadData          `json:"result,omitempty"`
		Errors       map[string]*errorData `json:"errors,omitempty"`
		Requirements []*requirementData    `json:"requirements,omitempty"`
	}

	payloadData struct {
		Type      *openapi.Schema `json:"type"`
		Example   interface{}     `json:"example,omitempty"`
		Streaming bool            `json:"streaming,omitempty"`
	}

	requirementData struct {
		Schemes []*schemeData `json:"schemes"`
		Scopes  []string      `json:"scopes"`
	}

	schemeData struct {
		Type        string      `json:"type"`
		Description string      `json:"description,omitempty"`
		Name        string      `json:"name"`
		In          string      `json:"in"`
		Scheme      string      `json:"scheme"`
		Flows       []*flowData `json:"flows,omitempty"`
	}

	errorData struct {
		Name        string          `json:"name"`
		Description string          `json:"description,omitempty,omitempty"`
		Type        *openapi.Schema `json:"type"`
		Temporary   bool            `json:"temporary,omitempty"`
		Timeout     bool            `json:"timeout,omitempty"`
		Fault       bool            `json:"fault,omitempty"`
	}

	flowData struct {
		Kind             string `json:"kind"`
		AuthorizationURL string `json:"authorizationURL,omitempty"`
		TokenURL         string `json:"tokenURL,omitempty"`
		RefreshURL       string `json:"refreshURL,omitempty"`
	}
)
