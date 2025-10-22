package design

import (
	. "goa.design/goa/v3/dsl"
	. "goa.design/plugins/v3/tools/dsl"
)

var _ = API("tool_example", func() {
	Title("Tool Example API")
	Description("Demonstrates the tools plugin with a single tool definition")
})

var _ = Service("inventory", func() {
	Description("Inventory service exposing tools for asset lookups")

	ToolSet("inventory_tools", func() {
		Tool("lookup_item", func() {
			Description("Retrieve an item from inventory")
			Payload(func() {
				Description("Item lookup parameters")
				Attribute("sku", String, "Inventory SKU", func() {
					MinLength(1)
				})
				Required("sku")
			})
			Result(func() {
				Description("Lookup result with optional details")
				Attribute("found", Boolean, "True when the item exists")
				Attribute("description", String, "Optional item description")
				Required("found")
			})
		})

		Tool("list_recent_items", func() {
			Description("List recently added inventory items")
			Payload(func() {
				Description("Filtering configuration for recent items")
				Attribute("limit", Int, "Maximum number of items to return", func() {
					Minimum(1)
					Maximum(100)
				})
				Required("limit")
			})
			Result(func() {
				Description("Recent items response envelope")
				Attribute("items", ArrayOf(String), "Item identifiers returned")
				Required("items")
			})
		})
	})
})
