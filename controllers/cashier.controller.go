package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
	"gorm.io/gorm"
)

// create cashier
func CreateCashierController(c *fiber.Ctx) error {
	db := database.Database()
	record := new(models.Cashier)
	if err := c.BodyParser(&record); err != nil {
		log.Fatal("failed to get data from the body")
	}
	db.Create(record)
	context := fiber.Map{"data": "", "msg": "cashier created successfully"}
	c.Status(http.StatusOK).JSON(context)
	return nil
}

// get all cashier list
func GetAllCashierController(c *fiber.Ctx) error {
	db := database.Database()

	var records []models.Cashier

	if err := db.Find(&records).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(records) <= 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "no data found"})
	}

	c.Status(http.StatusOK).JSON(fiber.Map{"msg": "data found", "data": records})
	return nil

}

// get cashier profile
func GetCashierProfileController(c *fiber.Ctx) error {
	db := database.Database()
	records := new(models.Cashier)
	id := c.Params("id")

	if err := db.First(&records, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"msg": "no data found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	c.Status(http.StatusOK).JSON(fiber.Map{"data": records, "msg": "profile retrieved successfully"})
	return nil
}

// update cashier profile
func UpdateCashierProfileController(c *fiber.Ctx) error {
	db := database.Database()
	id := c.Params("id")
	record := new(models.Cashier)
	db.First(&record, id)
	context := fiber.Map{"msg": ""}
	if record.ID == 0 {
		context["msg"] = "no data found"
		c.Status(http.StatusNotFound).JSON(context)
	}
	if err := c.BodyParser(&record); err != nil {
		context["msg"] = "failed to parse request body"
		return c.Status(http.StatusBadRequest).JSON(context)
	}
	result := db.Save(record)
	if result.Error != nil {
		context["msg"] = "an error occured while updating the data"
		c.Status(http.StatusInternalServerError).JSON(context)
	}
	context["msg"] = "data updated successfully"
	c.Status(http.StatusCreated).JSON(context)
	return nil
}

// delete cashier profile
func DeleteCashierProfileController(c *fiber.Ctx) error {
	db := database.Database()
	context := fiber.Map{"msg": ""}
	params := c.Params("id")
	var record models.Cashier
	db.First(&record, params)
	if record.ID == 0 {
		fmt.Println("no data found")
		context["msg"] = "Data Not Found!"
		return c.JSON(context)
	}
	db.Delete(&record)
	context["msg"] = "deleted successful"
	c.Status(200).JSON(context)
	return nil
}
