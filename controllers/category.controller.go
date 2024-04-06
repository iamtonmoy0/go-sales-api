package controllers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
	"gorm.io/gorm"
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
func GetAllCategoryController(c *fiber.Ctx) error {
	db := database.Database()
	records := []*models.Category{}
	// query
	query := db.Find(&records)
	if query.Error != nil {
		log.Fatal("Failed to fetch categories: ", query.Error)
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}
	c.Status(http.StatusOK).JSON(fiber.Map{
		"data": records,
	})
	return nil
}

// get single category
func GetSingleCategoryController(c *fiber.Ctx) error {
	db := database.Database()
	id := c.Params("id")
	var record models.Category
	result := db.First(&record, id)
	if result.Error == gorm.ErrRecordNotFound {
		c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "category not found",
		})
	}
	c.Status(http.StatusOK).JSON(fiber.Map{"msg": ""})
	return nil
}

// update category
func UpdateCategoryController(c *fiber.Ctx) error {
	db := database.Database()
	var record models.Category
	id := c.Params("id")
	err := c.BodyParser(&record)
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payload",
		})
	}
	//  check if the user exists in the database
	res := db.Model(&models.Category{}).Where("id = ?", id).Updates(&record)
	if res.Error != nil {
		c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "failed to update category",
		})
	}
	c.Status(http.StatusOK).JSON(fiber.Map{
		"msg": "successfully updated the category",
	})
	return nil
}

// delete category
func DeleteCategoryController(c *fiber.Ctx) error {
	db := database.Database()
	id := c.Params("id")
	record := &models.Category{}
	result := db.Delete(record, id)
	if result.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}
	c.Status(http.StatusOK).JSON(fiber.Map{"msg": "Successfully deleted Category!"})
	return nil
}
