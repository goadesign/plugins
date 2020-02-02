package testdata

import (
	. "goa.design/goa/dsl"
	cors "goa.design/plugins/cors/dsl"
)

var SimpleOriginDSL = func() {
	Service("SimpleOrigin", func() {
		cors.Origin("SimpleOrigin")
		Method("SimpleOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var RegexpOriginDSL = func() {
	Service("RegexpOrigin", func() {
		cors.Origin("/.*RegexpOrigin.*/")
		Method("RegexpOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var MultiOriginDSL = func() {
	Service("MultiOrigin", func() {
		cors.Origin("MultiOrigin1", func() {
			cors.Headers("X-Shared-Secret")
			cors.Methods("GET", "POST")
			cors.Expose("X-Time")
			cors.MaxAge(600)
			cors.Credentials()
		})
		cors.Origin("/.*MultiOrigin2.*/", func() {
			cors.Methods("GET", "POST")
			cors.Expose("X-Time", "X-Api-Version")
			cors.MaxAge(100)
		})
		Method("MultiOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var OriginFileServerDSL = func() {
	Service("OriginFileServer", func() {
		cors.Origin("OriginFileServer")
		Files("/file.json", "./file.json")
	})
}

var OriginMultiEndpointDSL = func() {
	Service("OriginMultiEndpoint", func() {
		cors.Origin("OriginMultiEndpoint")
		Method("OriginMultiEndpointGet", func() {
			HTTP(func() {
				GET("/{:id}")
			})
		})
		Method("OriginMultiEndpointPost", func() {
			HTTP(func() {
				POST("/")
			})
		})
		Method("OriginMultiEndpointOptions", func() {
			HTTP(func() {
				OPTIONS("/ids/{:id}")
			})
		})
	})
}

var MultiServiceSameOriginDSL = func() {
	Service("FirstService", func() {
		cors.Origin("SimpleOrigin")
		Method("SimpleOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
	Service("SecondService", func() {
		cors.Origin("SimpleOrigin")
		Method("SimpleOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}
