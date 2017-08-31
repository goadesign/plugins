package design

// Root is the design root expression.
var Root = new(RootExpr)

type (
	// RootExpr keeps track of the security schemes defined in the design
	RootExpr struct {
		// Schemes list the registered security schemes.
		Schemes []*SchemeExpr
		// APISecurity list the API level security requirements.
		APISecurity []*SecurityExpr
		// ServiceSecurity list service level security requirements.
		ServiceSecurity []*ServiceSecurityExpr
		// EndpointSecurity list endpoint level security requirements.
		EndpointSecurity []*EndpointSecurityExpr
		// FileServerSecurity list file server level security
		// requirements.
		FileServerSecurity []*FileServerSecurityExpr
	}
)
