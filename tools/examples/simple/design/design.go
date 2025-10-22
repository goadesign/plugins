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

	Method("ReserveStock", func() {
		Description("Reserve inventory units for a pending order")
		Payload(func() {
			Description("Order reservation parameters")
			Attribute("sku", String, "Inventory SKU to reserve", func() {
				MinLength(1)
			})
			Attribute("quantity", Int, "Number of units to reserve", func() {
				Minimum(1)
			})
			Required("sku", "quantity")
		})
		Result(func() {
			Description("Reservation outcome details")
			Attribute("reserved", Boolean, "True when the reservation succeeded")
			Attribute("reservation_id", String, "Identifier assigned to the reservation")
			Required("reserved")
		})
	})

	Method("SyncWarehouse", func() {
		Description("Synchronize stock levels with an external warehouse system")
		Payload(func() {
			Description("Warehouse synchronization payload")
			Attribute("warehouse_id", String, "External warehouse identifier", func() {
				MinLength(1)
			})
			Attribute("items", MapOf(String, Int), "Per-item stock counts to reconcile")
			Required("warehouse_id", "items")
			Meta("struct:pkg:path", "inventory/syncpayload")
		})
		Result(func() {
			Description("Synchronization response summary")
			Attribute("accepted", Boolean, "True when the update is accepted")
			Attribute("errors", ArrayOf(String), "Optional per-item validation errors")
			Required("accepted")
			Meta("struct:pkg:path", "inventory/syncresult")
		})
	})

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

	ToolSet("inventory_method_tools", func() {
		ToolFromMethod("ReserveStock")
		ToolFromMethod("SyncWarehouse")
	})
})
