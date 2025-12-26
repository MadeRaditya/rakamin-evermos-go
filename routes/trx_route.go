package routes

import (
	"backend-evermos/controllers"
	"backend-evermos/middlewares"

	"github.com/gofiber/fiber/v2"
)

func TrxRoute(api fiber.Router) {
	trx := api.Group("/trx", middlewares.JWTAuth)

	trx.Get("/", controllers.GetAllTrx)
	trx.Get("/:id", controllers.GetTrxByID)
	trx.Post("/", controllers.CreateTrx)
}
