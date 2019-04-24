package testdata

import (
	. "goa.design/goa/dsl"
	_ "goa.design/plugins/docs"
)

var APIOnly = func() {
	API("API", func() {
		Server("Host1", func() {
			Services("Service")
			Host("dev", func() {
				URI("http://example:8090")
			})
		})
		Server("Host2", func() {
			Services("Service")
			Host("dev", func() {
				URI("http://example:8090")
			})
		})
	})
}

var NoPayloadNoReturn = func() {
	API("SingleService", func() {
		Server("SingleHost", func() {
			Services("Service")
			Host("dev", func() {
				URI("http://example:8090")
			})
		})
	})
	Service("Service", func() {
		Method("Method", func() {
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var PrimitivePayloadNoReturn = func() {
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Payload(String)
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var ArrayPayloadNoReturn = func() {
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Payload(ArrayOf(String))
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var MapPayloadNoReturn = func() {
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Payload(MapOf(String, Int32))
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var UserPayloadNoReturn = func() {
	var User = Type("User", func() {
		Attribute("att1", String)
		Attribute("att2", Int)
	})
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Payload(User)
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var NoPayloadPrimitiveReturn = func() {
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Result(String)
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var NoPayloadArrayReturn = func() {
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Result(ArrayOf(String))
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var NoPayloadMapReturn = func() {
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Result(MapOf(String, Int32))
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}

var NoPayloadUserReturn = func() {
	var User = Type("User", func() {
		Attribute("att1", String)
		Attribute("att2", Int)
	})
	API("Test API", func() {})
	Service("Service", func() {
		Method("Method", func() {
			Result(User)
			HTTP(func() {
				GET("/")
			})
			GRPC(func() {})
		})
	})
}
