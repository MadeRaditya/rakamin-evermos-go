package routes

import (
	"github.com/MadeRaditya/rakamin-evermos-go/controllers"
	"github.com/MadeRaditya/rakamin-evermos-go/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(api fiber.Router) {
	product := api.Group("/product", middlewares.JWTAuth)

	product.Get("/", controllers.GetAllProduct)
	product.Post("/", controllers.CreateProduct)
	product.Get("/:id", controllers.GetProductByID)
	product.Put("/:id", controllers.UpdateProduct)
	product.Delete("/:id", controllers.DeleteProduct)
}
