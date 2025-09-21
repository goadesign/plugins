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
