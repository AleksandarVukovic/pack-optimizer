package design

import . "goa.design/goa/v3/dsl"

var _ = API("pack-optimizer", func() {
	Title("Pack Optimizer")
	Description("Service to calculate optimal pack combinations for a given number of items")
	Server("pack-optimizer", func() {
		Host("localhost", func() {
			URI("http://localhost:8080")
		})
	})
})

var _ = Service("optimizer", func() {
	Description("Service calculates optimal pack combinations for a given number of items.")

	Method("health", func() {
		Description("Returns the health status of the service.")

		Result(func() {
			Field(1, "status", String, "Health status of the service", func() {
				Example("ok")
			})
			Required("status")
		})

		Error("internal_server_error", ErrorResult)

		HTTP(func() {
			GET("/health")
			Response(StatusOK)
			Response("internal_server_error", StatusInternalServerError)
		})
	})

	Method("getPackSizes", func() {
		Description("Get the current pack sizes.")
		Result(func() {
			Field(1, "sizes", ArrayOf(Int), "Current pack sizes", func() {
				Example([]int{250, 500, 1000})
			})
			Required("sizes")
		})

		Error("internal_server_error", ErrorResult)

		HTTP(func() {
			GET("/packs/sizes")
			Response(StatusOK)
			Response("internal_server_error", StatusInternalServerError)
		})
	})

	Method("updatePackSizes", func() {
		Description("Update the pack sizes.")
		Payload(func() {
			Field(1, "sizes", ArrayOf(Int), "New pack sizes to update", func() {
				Example([]int{250, 500, 1000})
			})
			Required("sizes")
		})

		Error("bad_request", ErrorResult)
		Error("internal_server_error", ErrorResult)

		HTTP(func() {
			PUT("/packs/sizes")
			Response(StatusNoContent)
			Response("bad_request", StatusBadRequest)
			Response("internal_server_error", StatusInternalServerError)
		})
	})

	Method("calculate", func() {
		Description("Calculate optimal pack combinations for a given number of items.")
		Payload(func() {
			Field(1, "quantity", Int, "Total number of items to pack", func() {
				Minimum(1)
				Example(500)
			})
			Required("quantity")
		})
		Result(func() {
			Field(1, "packs", ArrayOf(Pack), "Optimal pack combinations (pack size -> quantity)")
			Required("packs")
		})

		Error("bad_request", ErrorResult)
		Error("internal_server_error", ErrorResult)

		HTTP(func() {
			GET("/packs/calculate")
			Param("quantity")
			Response(StatusOK)
			Response("bad_request", StatusBadRequest)
			Response("internal_server_error", StatusInternalServerError)
		})
	})

})

var Pack = Type("Pack", func() {
	Description("Pack represents a pack size and the quantity needed.")
	Field(1, "size", Int, "Size of the pack")
	Field(2, "quantity", Int, "Quantity of this pack size needed")
	Required("size", "quantity")
})
