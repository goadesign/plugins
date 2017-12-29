package design

import (
	. "goa.design/goa/http/design"
	. "goa.design/plugins/cors/dsl"
)

var _ = Service("sommelier", func() {
	Description("The sommelier service retrieves bottles given a set of criteria.")
	HTTP(func() {
		Path("/sommelier")
	})
	Origin("/.*localhost*/", func() {
		Methods("GET", "POST")
		Expose("X-Time")
		MaxAge(600)
	})
	Method("pick", func() {
		Payload(Criteria)
		Result(CollectionOf(StoredBottle, func() {
			View("default")
		}))
		Error("no_criteria", String, "Missing criteria")
		Error("no_match", String, "No bottle matched given criteria")
		HTTP(func() {
			POST("/")
			Response(StatusOK)
			Response("no_criteria", StatusBadRequest)
			Response("no_match", StatusNotFound)
		})
	})
})
