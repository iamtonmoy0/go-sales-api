package controllers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
)

// create product
func CreateProductController(c *fiber.Ctx) error {
	db := database.Database()
	var record models.Product
	context := fiber.Map{"msg": ""}
	if err := c.BodyParser(&record); err != nil {
		log.Fatalf("failed to parse the data", err)
		context["msg"] = "failed to parse the body"
		c.Status(http.StatusInternalServerError).JSON(context)
	}
	err := db.Create(&record).Error
	if err != nil {
		log.Println(err)
		context["msg"] = "failed to create product"
		c.Status(http.StatusBadRequest).JSON(context)
	}
	context["msg"] = "product created!"
	return c.Status(http.StatusCreated).JSON(context)
	return nil
}

// get all product
func GetAllProductController(c *fiber.Ctx) error {
	db := database.Database()
	var records []models.Product
	result := db.Find(&records)
	if result.Error != nil {
		log.Println(result.Error)
		return c.Status(http.StatusNotFound).SendString("No products found!")
	}
	c.Status(http.StatusOK).JSON(result)
	return nil
}

// get single product
func GetSingleProductController(c *fiber.Ctx) error {
	db := database.Database()
	id := c.Params("id")
	var record models.Product
	result := db.First(&record, id)
	if result.Error != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "failed to get the product"})
		return nil
	}
	c.Status(http.StatusOK).JSON(result)
	return nil
}

// update product

func UpdateProductController(c *fiber.Ctx) error {
	db := database.Database()
	id := c.Params("id")
	var oldRecord models.Product
	var newRecord models.Product
	if err := c.BodyParser(&newRecord); err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "failed to parse the data"})
		return err
	}
	if err := db.Where("id = ?", id).First(&oldRecord).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": "failed to find product with given id",
			"err": err.Error(),
		})
		return err
	}
	if err := db.Model(&oldRecord).Updates(&newRecord).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": "product not updated",
			"err": err.Error(),
		})
		return err
	}
	c.Status(http.StatusOK).JSON(fiber.Map{
		"msg": "successfully updated the product",
	})
	return nil
}

// delete product
func DeleteProductController(c *fiber.Ctx) error {
	db := database.Database()
	id := c.Params("id")
	var record models.Product
	if err := db.First(&record, id).Delete(&record).Error; err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Record not found!"})
		return nil
	}
	c.Status(http.StatusOK).SendString("The product has been deleted.")
	return nil
}
