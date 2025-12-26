package testdata

import (
    . "goa.design/goa/v3/dsl"
)

// Single service with HTTP and one server.
var SingleServiceDSL = func() {
    API("calc", func() {
        Server("edge", func() {
            Services("calc")
            Host("dev", func() { URI("http://localhost:8080") })
        })
    })
    Service("calc", func() {
        Method("add", func() {
            Payload(func() { Attribute("a", Int, "a"); Attribute("b", Int, "b"); Required("a", "b") })
            Result(Int)
            HTTP(func() { GET("/add") })
        })
    })
}

// Two services under one server.
var MultiServiceDSL = func() {
    API("multi", func() {
        Server("api", func() {
            Services("s1", "s2")
            Host("dev", func() { URI("http://localhost:8080") })
        })
    })
    Service("s1", func() {
        Method("m1", func() {
            Result(String)
            HTTP(func() { GET("/s1") })
        })
    })
    Service("s2", func() {
        Method("m2", func() {
            Result(String)
            HTTP(func() { GET("/s2") })
        })
    })
}

// WebSocket streaming service (bidirectional): should trigger WS import and upgrader usage.
var WebSocketServiceDSL = func() {
    API("chatapi", func() {
        Server("chat", func() {
            Services("chat")
            Host("dev", func() { URI("http://localhost:8080") })
        })
    })
    var Message = Type("Message", func() {
        Attribute("text", String)
        Required("text")
    })
    Service("chat", func() {
        Method("room", func() {
            StreamingPayload(Message)
            StreamingResult(Message)
            HTTP(func() { GET("/chat") })
        })
    })
}

// Single service with only Files() endpoints; used to validate mains adds the
// correct number of http.FileSystem arguments to the HTTP server constructor.
var FileServerServiceDSL = func() {
    API("fsapi", func() {
        Server("edge", func() {
            Services("static")
            Host("dev", func() { URI("http://localhost:8080") })
        })
    })
    Service("static", func() {
        HTTP(func() { Path("/") })
        Files("/f1.json", "/assets/f1.json")
        Files("/f2.json", "/assets/f2.json")
        Files("/f3.json", "/assets/f3.json")
    })
}
