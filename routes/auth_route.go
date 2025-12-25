package routes

import (
	"backend-evermos/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)
}
