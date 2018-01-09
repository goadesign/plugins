package testdata

import (
	. "goa.design/goa/design"
	. "goa.design/plugins/cors/dsl"
)

var SimpleOriginDSL = func() {
	Service("SimpleOrigin", func() {
		Origin("SimpleOrigin")
		Method("SimpleOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var RegexpOriginDSL = func() {
	Service("RegexpOrigin", func() {
		Origin("/.*RegexpOrigin.*/")
		Method("RegexpOriginMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var MultiOriginDSL = func() {
	Service("MultiOrigin", func() {
		Origin("MultiOrigin1", func() {
			Headers("X-Shared-Secret")
			Methods("GET", "POST")
			Expose("X-Time")
			MaxAge(600)
			Credentials()
		})
		Origin("/.*MultiOrigin2.*/", func() {
			Methods("GET", "POST")
			Expose("X-Time", "X-Api-Version")
			MaxAge(100)
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
		Origin("OriginFileServer")
		Files("/file.json", "./file.json")
	})
}

var OriginMultiEndpointDSL = func() {
	Service("OriginMultiEndpoint", func() {
		Origin("OriginMultiEndpoint")
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
