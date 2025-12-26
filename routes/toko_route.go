package routes

import (
	"backend-evermos/controllers"
	"backend-evermos/middlewares"

	"github.com/gofiber/fiber/v2"
)

func TokoRoute(api fiber.Router) {
	tokoRoute := api.Group("/toko", middlewares.JWTAuth)

	tokoRoute.Get("/", controllers.GetAllToko)
	tokoRoute.Get("/my", controllers.GetMyToko)
	tokoRoute.Get("/:id", controllers.GetTokoByID)
	tokoRoute.Put("/:id_toko", controllers.UpdateTokoByID)
}
