package testdata

import (
	. "goa.design/goa/design"
	"goa.design/plugins/security/design"
	. "goa.design/plugins/security/dsl"
)

// TopLevelSchemes contains all defined top-level security schemes
// which are added to a security RootExpression.
var TopLevelSchemes = []*design.SchemeExpr{BasicAuth, APIKeyAuth, OAuth2Auth, JWTAuth}

var BasicAuth = BasicAuthSecurity("basic")

var APIKeyAuth = APIKeySecurity("api_key")

var OAuth2Auth = OAuth2Security("pass", func() {
	PasswordFlow("/token", "/refresh")
	Scope("api:write", "Write acess")
	Scope("api:read", "Read access")
})

var JWTAuth = JWTSecurity("jwt", func() {
	Scope("api:read", "Read-only access")
	Scope("api:write", "Read and write access")
})

var BasicAuthDSL = func() {
	Service("BasicAuth", func() {
		Method("login", func() {
			Security(BasicAuth)
			Payload(func() {
				Username("user", String)
				Password("pass", String)
			})
			HTTP(func() {
				POST("/login")
			})
		})
	})
}

var OAuth2DSL = func() {
	Service("OAuth2", func() {
		Method("login", func() {
			Security(OAuth2Auth)
			Payload(func() {
				AccessToken("token", String)
			})
			HTTP(func() {
				POST("/login")
			})
		})
	})
}

var OAuth2InParamDSL = func() {
	Service("OAuth2InParam", func() {
		Method("login", func() {
			Security(OAuth2Auth)
			Payload(func() {
				AccessToken("token", String)
			})
			HTTP(func() {
				POST("/login")
				Param("token:t")
			})
		})
	})
}

var JWTDSL = func() {
	Service("JWT", func() {
		Method("login", func() {
			Security(JWTAuth)
			Payload(func() {
				Token("token", String)
			})
			HTTP(func() {
				POST("/login")
				Header("Authorization")
			})
		})
	})
}

var APIKeyDSL = func() {
	Service("APIKey", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/login")
				Header("key:Authorization")
			})
		})
	})
}

var APIKeyInParamDSL = func() {
	Service("APIKeyInParam", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/login")
				Param("key")
			})
		})
	})
}

var MultipleAndDSL = func() {
	Service("MultipleAnd", func() {
		Method("login", func() {
			Security(BasicAuth, APIKeyAuth)
			Payload(func() {
				Username("user", String)
				Password("password", String)
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/login")
				Param("key:k")
			})
		})
	})
}

var MultipleOrDSL = func() {
	Service("MultipleOr", func() {
		Method("login", func() {
			Security(BasicAuth)
			Security(APIKeyAuth)
			Payload(func() {
				Username("user", String)
				Password("password", String)
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/login")
				Param("key:k")
			})
		})
	})
}
