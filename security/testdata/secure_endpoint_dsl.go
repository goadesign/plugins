package testdata

import (
	. "goa.design/goa/design"
	. "goa.design/plugins/security/dsl"
)

var BasicAuth = BasicAuthSecurity("basic")

var JWTAuth = JWTSecurity("jwt", func() {
	Scope("api:read", "Read-only access")
	Scope("api:write", "Read and write access")
})

var APIKeyAuth = APIKeySecurity("api_key")

var OAuth2Implicit = OAuth2Security("implicit", func() {
	ImplicitFlow("/authorization", "/refresh")
	Scope("api:write", "Write acess")
	Scope("api:read", "Read access")
})

var OAuth2AuthorizationCode = OAuth2Security("authCode", func() {
	AuthorizationCodeFlow("/authorization", "/token", "/refresh")
	Scope("api:write", "Write acess")
	Scope("api:read", "Read access")
})

var OAuth2Password = OAuth2Security("pass", func() {
	PasswordFlow("/token", "/refresh")
	Scope("api:write", "Write acess")
	Scope("api:read", "Read access")
})

var OAuth2ClientCredentials = OAuth2Security("clicred", func() {
	ClientCredentialsFlow("/token", "/refresh")
	Scope("api:write", "Write acess")
	Scope("api:read", "Read access")
})

var EndpointWithoutRequirementDSL = func() {
	Service("EndpointWithoutRequirement", func() {
		Method("Unsecure", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var EndpointNoSecurityDSL = func() {
	Service("EndpointNoSecurity", func() {
		Security(BasicAuth)
		Method("NoSecurity", func() {
			NoSecurity()
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var EndpointsWithServiceRequirementsDSL = func() {
	Service("EndpointsWithServiceRequirements", func() {
		Security(BasicAuth)
		Method("SecureWithRequirements", func() {
			HTTP(func() {
				GET("/")
			})
		})
		Method("AlsoSecureWithRequirements", func() {
			HTTP(func() {
				POST("/")
			})
		})
	})
}

var EndpointsWithRequirementsDSL = func() {
	Service("EndpointsWithRequirements", func() {
		Method("SecureWithRequirements", func() {
			Security(BasicAuth)
			HTTP(func() {
				GET("/")
			})
		})
		Method("DoublySecureWithRequirements", func() {
			Security(BasicAuth, JWTAuth)
			HTTP(func() {
				POST("/")
			})
		})
	})
}

var EndpointWithRequiredScopesDSL = func() {
	Service("EndpointWithRequiredScopes", func() {
		Method("SecureWithRequiredScopes", func() {
			Security(JWTAuth, func() {
				Scope("api:read")
				Scope("api:write")
			})
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var EndpointWithAPIKeyOverrideDSL = func() {
	Service("EndpointWithAPIKeyOverride", func() {
		Security(BasicAuth)
		Method("SecureWithAPIKeyOverride", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var EndpointWithOAuth2DSL = func() {
	Service("EndpointWithOAuth2", func() {
		Method("SecureWithOAuth2", func() {
			Security(OAuth2Implicit, OAuth2AuthorizationCode, OAuth2Password, OAuth2ClientCredentials)
			Payload(func() {
				AccessToken("token", String)
			})
			HTTP(func() {
				GET("/")
			})
		})
	})
}
