package routes

import (
	"github.com/MadeRaditya/rakamin-evermos-go/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProvCityRoute(api fiber.Router) {
	prov := api.Group("/provcity")

	prov.Get("/listprovincies", controllers.GetListProvinces)
	prov.Get("/listcities/:prov_id", controllers.GetListCities)
	prov.Get("/detailprovince/:prov_id", controllers.GetDetailProvince)
	prov.Get("/detailcity/:city_id", controllers.GetDetailCity)
}
