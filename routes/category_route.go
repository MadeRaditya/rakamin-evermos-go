package routes

import (
	"github.com/MadeRaditya/rakamin-evermos-go/controllers"
	"github.com/MadeRaditya/rakamin-evermos-go/middlewares"

	"github.com/gofiber/fiber/v2"
)

func CategoryRoute(api fiber.Router) {
	category := api.Group("/category", middlewares.JWTAuth)

	category.Get("/", controllers.GetAllCategory)
	category.Post("/", controllers.CreateCategory)
	category.Get("/:id", controllers.GetCategoryByID)
	category.Put("/:id", controllers.UpdateCategory)
	category.Delete("/:id", controllers.DeleteCategory)
}
