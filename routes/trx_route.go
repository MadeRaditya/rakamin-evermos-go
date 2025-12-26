package routes

import (
	"github.com/MadeRaditya/rakamin-evermos-go/controllers"
	"github.com/MadeRaditya/rakamin-evermos-go/middlewares"

	"github.com/gofiber/fiber/v2"
)

func TrxRoute(api fiber.Router) {
	trx := api.Group("/trx", middlewares.JWTAuth)

	trx.Get("/", controllers.GetAllTrx)
	trx.Get("/:id", controllers.GetTrxByID)
	trx.Post("/", controllers.CreateTrx)
}
