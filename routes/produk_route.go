package routes

import (
	"backend-evermos/controllers"
	"backend-evermos/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ProductRouter(api fiber.Router) {
	productRoute := api.Group("/product", middlewares.JWTAuth)

	productRoute.Get("/", controllers.GetAllProduk)
	productRoute.Post("/", controllers.CreateProduct)

}
