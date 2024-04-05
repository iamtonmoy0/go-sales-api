package controllers

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
	"github.com/iamtonmoy0/go-sales-api/utils"
)

type cashier struct {
	Email    string `json:"email" gorm:"not null"`
	Password string `json:"password gorm:"not null"`
}

// login controller
func Login(c *fiber.Ctx) error {
	db := database.Database()
	var record cashier
	if err := c.BodyParser(&record); err != nil {
		log.Fatal("failed to parse the body", err)
	}
	var cashier models.Cashier
	if err := db.Where("email=?", record.Email).First(&cashier); err != nil {
		log.Fatal("No user found", err)
		c.Status(http.StatusBadGateway).JSON(fiber.Map{"msg": "Invalid user credentials"})
	}
	// checking password
	isValid := utils.CheckPasswordHash(record.Password, cashier.Password)
	if !isValid {
		log.Println("Failed to check password: ")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"msg": "Invalid user credentials",
		})
	}
	// generating token
	jwtToken, err := utils.GenerateToken(cashier.ID)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	c.Status(http.StatusOK).JSON(fiber.Map{
		"token": jwtToken,
	})
	return nil
}

// logout controller
func Logout(c *fiber.Ctx) error {
	return nil
}
