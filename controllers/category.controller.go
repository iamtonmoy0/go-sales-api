package controllers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
)

// create category
func CreateCategoryController(c *fiber.Ctx) error {
	db := database.Database()
	var record models.Category
	if err := c.BodyParser(&record); err != nil {
		log.Fatal(err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "failed to parse body"})
	}
	if err := db.Create(&record); err != nil {
		log.Fatal("failed to create category", err)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "failed to create category "})
	}
	c.Status(http.StatusOK).JSON(fiber.Map{"msg": "new category created"})
	return nil

}

// get all category
func GetAllCategoryController(c *fiber.Ctx) error { return nil }

// get single category
func GetSingleCategoryController(c *fiber.Ctx) error { return nil }

// update category
func UpdateCategoryController(c *fiber.Ctx) error { return nil }

// delete category
func DeleteCategoryController(c *fiber.Ctx) error { return nil }
