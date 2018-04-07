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

var BasicAuthRequiredDSL = func() {
	Service("BasicAuthRequired", func() {
		Method("login", func() {
			Security(BasicAuth)
			Payload(func() {
				Username("user", String)
				Password("pass", String)
				Attribute("id", String)
				Required("user", "pass")
			})
			HTTP(func() {
				POST("/login")
				Param("id")
			})
		})
	})
}

var OAuth2DSL = func() {
	Service("OAuth2", func() {
		Method("login", func() {
			Security(OAuth2Auth, func() {
				Scope("api:write")
			})
			Payload(func() {
				AccessToken("token", String)
			})
			HTTP(func() {
				POST("/login")
			})
		})
	})
}

var OAuth2RequiredDSL = func() {
	Service("OAuth2Required", func() {
		Method("login", func() {
			Security(OAuth2Auth, func() {
				Scope("api:write")
			})
			Payload(func() {
				AccessToken("token", String)
				Required("token")
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

var OAuth2InParamRequiredDSL = func() {
	Service("OAuth2InParamRequired", func() {
		Method("login", func() {
			Security(OAuth2Auth)
			Payload(func() {
				AccessToken("token", String)
				Required("token")
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
			Security(JWTAuth, func() {
				Scope("api:read")
			})
			Payload(func() {
				Token("token", String)
			})
			HTTP(func() {
				POST("/login")
				Header("token:Authorization")
			})
		})
	})
}

var JWTRequiredDSL = func() {
	Service("JWTRequired", func() {
		Method("login", func() {
			Security(JWTAuth, func() {
				Scope("api:read")
			})
			Payload(func() {
				Token("token", String)
				Required("token")
			})
			HTTP(func() {
				POST("/login")
				Header("token:Authorization")
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

var APIKeyRequiredDSL = func() {
	Service("APIKeyRequired", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
				Required("key")
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

var APIKeyInParamRequiredDSL = func() {
	Service("APIKeyInParamRequired", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
				Required("key")
			})
			HTTP(func() {
				POST("/login")
				Param("key")
			})
		})
	})
}

var APIKeyInBodyDSL = func() {
	Service("APIKeyInBody", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/login")
				Body("key")
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

var MultipleAndRequiredDSL = func() {
	Service("MultipleAndRequired", func() {
		Method("login", func() {
			Security(BasicAuth, APIKeyAuth)
			Payload(func() {
				Username("user", String)
				Password("password", String)
				APIKey("api_key", "key", String)
				Required("user", "password", "key")
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

var MultipleOrRequiredDSL = func() {
	Service("MultipleOrRequired", func() {
		Method("login", func() {
			Security(BasicAuth)
			Security(APIKeyAuth)
			Payload(func() {
				Username("user", String)
				Password("password", String)
				APIKey("api_key", "key", String)
				Required("user", "password", "key")
			})
			HTTP(func() {
				POST("/login")
				Param("key:k")
			})
		})
	})
}

var MultipleSchemesWithParamsDSL = func() {
	Service("MultipleSchemesWithParams", func() {
		Method("login", func() {
			Security(JWTAuth)
			Security(APIKeyAuth)
			Payload(func() {
				Token("token", String)
				APIKey("api_key", "key", String)
				Attribute("id", Int)
				Attribute("user_agent", String)
				Attribute("name", String)
				Attribute("description", String)
				Required("token")
			})
			HTTP(func() {
				POST("/login/{id}")
				Param("key:k")
				Param("name")
				Header("user_agent:User-Agent")
			})
		})
	})
}

var SchemesInTypeDSL = func() {
	var Schemes = Type("Schemes", func() {
		APIKey("api_key", "key", String)
		AccessToken("token", String)
	})
	Service("SchemesInTypeDSL", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Security(OAuth2Auth)
			Payload(Schemes)
			HTTP(func() {
				POST("/login")
				Param("key:k")
			})
		})
	})
}

var SchemesInTypeRequiredDSL = func() {
	var Schemes = Type("SchemesRequired", func() {
		APIKey("api_key", "key", String)
		AccessToken("token", String)
		Required("key", "token")
	})
	Service("SchemesInTypeRequiredDSL", func() {
		Method("login", func() {
			Security(APIKeyAuth)
			Security(OAuth2Auth)
			Payload(Schemes)
			HTTP(func() {
				POST("/login")
				Param("key:k")
			})
		})
	})
}

var SameSchemeMultipleEndpoints = func() {
	Service("SameSchemeMultipleEndpoints", func() {
		Method("method1", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/method1")
				Param("key:k")
			})
		})
		Method("method2", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				POST("/method1")
				Header("key:Authorization")
			})
		})
	})
}

var SingleServiceDSL = func() {
	Service("SingleService", func() {
		Method("Method", func() {
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

var MultipleServicesDSL = func() {
	Service("ServiceWithAPIKeyAuth", func() {
		Method("Method", func() {
			Security(APIKeyAuth)
			Payload(func() {
				APIKey("api_key", "key", String)
			})
			HTTP(func() {
				GET("/")
			})
		})
	})
	Service("ServiceWithJWTAndBasicAuth", func() {
		Security(BasicAuth, JWTAuth)
		Method("Method", func() {
			Payload(func() {
				Username("user", String)
				Password("pass", String)
				Token("token", String)
			})
			HTTP(func() {
				GET("/")
			})
		})
	})
	Service("ServiceWithNoSecurity", func() {
		Method("Method", func() {
			Payload(func() {
				Attribute("a", String)
			})
			HTTP(func() {
				GET("/{a}")
			})
		})
	})
}
