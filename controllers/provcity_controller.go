package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Province struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         string `json:"id"`
	ProvinceID string `json:"province_id"`
	Name       string `json:"name"`
}

const BaseURL = "https://www.emsifa.com/api-wilayah-indonesia/api"

func featchAPI(url string, target interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func GetListProvinces(c *fiber.Ctx) error {
	url := fmt.Sprintf("%s/provinces.json", BaseURL)

	var provinces []Province
	if err := featchAPI(url, &provinces); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to get data from external API",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    provinces,
	})
}

func GetListCities(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	url := fmt.Sprintf("%s/regencies/%s.json", BaseURL, provID)

	var cities []City
	if err := featchAPI(url, &cities); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to get data from external API",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    cities,
	})
}

func GetDetailProvince(c *fiber.Ctx) error {
	provID := c.Params("prov_id")

	url := fmt.Sprintf("%s/provinces.json", BaseURL)

	var provinces []Province
	if err := featchAPI(url, &provinces); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status": false, "message": "Failed to external API", "errors": []string{err.Error()}, "data": nil,
		})
	}

	var foundProv *Province
	for _, p := range provinces {
		if p.ID == provID {
			foundProv = &p
			break
		}
	}

	if foundProv == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false, "message": "Not Found", "errors": []string{"Province not found"}, "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    foundProv,
	})
}

func GetDetailCity(c *fiber.Ctx) error {
	cityID := c.Params("city_id")

	if len(cityID) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false, "message": "Invalid ID", "errors": []string{"City ID invalid"}, "data": nil,
		})
	}

	provID := cityID[:2]

	url := fmt.Sprintf("%s/regencies/%s.json", BaseURL, provID)

	var cities []City
	if err := featchAPI(url, &cities); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"status": false, "message": "Failed to external API", "errors": []string{err.Error()}, "data": nil,
		})
	}

	var foundCity *City
	for _, city := range cities {
		if city.ID == cityID {
			foundCity = &city
			break
		}
	}

	if foundCity == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false, "message": "Not Found", "errors": []string{"City not found"}, "data": nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    foundCity,
	})
}
