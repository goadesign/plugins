package testdata

import (
	. "goa.design/goa/dsl"
)

var SimpleServiceDSL = func() {
	Service("SimpleService", func() {
		Method("SimpleMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var WithPayloadDSL = func() {
	Service("WithPayloadService", func() {
		Method("WithPayloadMethod", func() {
			Payload(func() {
				Attribute("id")
			})
			HTTP(func() {
				GET("/")
			})
		})
	})
}

var WithErrorDSL = func() {
	Service("WithErrorService", func() {
		Method("WithErrorMethod", func() {
			Error("bad_request")
			HTTP(func() {
				GET("/")
				Response("bad_request", StatusBadRequest)
			})
		})
	})
}

var MultiEndpointDSL = func() {
	Service("MultiEndpointService", func() {
		Error("bad_request")
		Method("Endpoint1", func() {
			Payload(func() {
				Attribute("id")
			})
			HTTP(func() {
				GET("/")
				Response("bad_request", StatusBadRequest)
			})
		})
		Method("Endpoint2", func() {
			HTTP(func() {
				POST("/")
				Response("bad_request", StatusBadRequest)
			})
		})
	})
}

var FileServerDSL = func() {
	Service("FileServerService", func() {
		Files("/1.json", "../file1.json")
		Files("/2.json", "../file2.json")
	})
}

var MixedDSL = func() {
	Service("MixedService", func() {
		Method("MixedMethod", func() {
			HTTP(func() {
				GET("/")
			})
		})
		Files("/1.json", "../mixed_file.json")
	})
}

var MultiServiceDSL = func() {
	Service("Service1", func() {
		Method("Method", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
	Service("Service2", func() {
		Method("Method", func() {
			HTTP(func() {
				GET("/")
			})
		})
	})
}
