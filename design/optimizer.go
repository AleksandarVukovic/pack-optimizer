package design

import . "goa.design/goa/v3/dsl"

var _ = Service("optimizer", func() {
	Description("The optimizer service calculates optimal pack combinations for a given number of items.")

	Method("getSizes", func() {
		Description("Get the current pack sizes.")
		Result(func() {
			Field(1, "sizes", ArrayOf(Int), "Current pack sizes")
			Required("sizes")
		})
		HTTP(func() {
			GET("/packs/sizes")
			Response(StatusOK)
			Response(StatusInternalServerError)
		})
	})

	Method("updateSizes", func() {
		Description("Update the pack sizes.")
		Payload(func() {
			Field(1, "sizes", ArrayOf(Int), "New pack sizes to update")
			Required("sizes")
		})
		HTTP(func() {
			PUT("/packs/sizes")
			Response(StatusNoContent)
			Response(StatusBadRequest)
			Response(StatusInternalServerError)
		})
	})

	Method("calculate", func() {
		Description("Calculate optimal pack combinations for a given number of items.")
		Payload(func() {
			Field(1, "totalItems", Int, "Total number of items to pack")
			Required("totalItems")
		})
		Result(func() {
			Field(1, "packs", ArrayOf(Pack), "Optimal pack combinations (pack size -> quantity)")
			Required("packs")
		})
		HTTP(func() {
			GET("/packs/calculate")
			Response(StatusOK)
			Response(StatusBadRequest)
			Response(StatusInternalServerError)
		})
	})

})

var Pack = Type("Pack", func() {
	Description("Pack represents a pack size and the quantity needed.")
	Field(1, "size", Int, "Size of the pack")
	Field(2, "quantity", Int, "Quantity of this pack size needed")
	Required("size", "quantity")
})
