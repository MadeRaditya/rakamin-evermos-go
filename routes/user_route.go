package routes

import (
	"backend-evermos/controllers/user"
	"backend-evermos/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(api fiber.Router) {
	userRoute := api.Group("/user", middlewares.JWTAuth)

	userRoute.Get("/", user.GetMyProfil)
	userRoute.Put("/", user.UpdateProfil)

	userRoute.Get("/alamat", user.GetAllAlamat)
	userRoute.Post("/alamat", user.CreateAlamat)
	userRoute.Get("/alamat/:id", user.GetAlamatByID)
	userRoute.Put("/alamat/:id", user.UpdateAlamat)
	userRoute.Delete("/alamat/:id", user.DeleteAlamat)
}
