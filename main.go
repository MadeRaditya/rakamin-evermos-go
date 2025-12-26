package main

import (
	"log"
	"os"

	"backend-evermos/database"
	"backend-evermos/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using system env")
	}

	database.Connect()

	app := fiber.New()
	app.Static("/public", "./public")
	api := app.Group("/api/v1")

	routes.AuthRoute(api)
	routes.CategoryRoute(api)
	routes.UserRoute(api)
	routes.TokoRoute(api)
	routes.ProductRouter(api)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	log.Fatal(app.Listen(":" + port))
}
