package routes

import (
	"github.com/MadeRaditya/rakamin-evermos-go/controllers"
	"github.com/MadeRaditya/rakamin-evermos-go/middlewares"

	"github.com/gofiber/fiber/v2"
)

func TokoRoute(api fiber.Router) {
	tokoRoute := api.Group("/toko", middlewares.JWTAuth)

	tokoRoute.Get("/", controllers.GetAllToko)
	tokoRoute.Get("/my", controllers.GetMyToko)
	tokoRoute.Get("/:id", controllers.GetTokoByID)
	tokoRoute.Put("/:id_toko", controllers.UpdateTokoByID)
}
