package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
	"github.com/iamtonmoy0/go-sales-api/utils"
	"gorm.io/gorm"
)

// login controller
func Login(c *fiber.Ctx) error {
	db := database.Database()
	var record models.Cashier

	// Parse request body
	if err := c.BodyParser(&record); err != nil {
		log.Fatal("Failed to parse the body", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"msg": "Failed to parse request body"})
	}

	// Query cashier by email
	var cashier models.Cashier
	if err := db.Where("email = ?", record.Email).First(&cashier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("No user found:", err)
			return c.Status(http.StatusBadGateway).JSON(fiber.Map{"msg": "Invalid user credentials"})
		}
		log.Fatal("Error fetching user:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "Failed to fetch user"})
	}

	// Check password
	isValid := utils.CheckPasswordHash(record.Password, cashier.Password)
	if !isValid {
		log.Println("Failed to check password")
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"msg": "Invalid user credentials"})
	}

	// Generate token
	jwtToken, err := utils.GenerateToken(cashier.ID)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "Failed to generate token"})
	}

	// Return token
	return c.Status(http.StatusOK).JSON(fiber.Map{"token": jwtToken})
}

// logout controller
func Logout(c *fiber.Ctx) error {
	// Clear the token from the client-side by setting an expired token
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
		HTTPOnly: true,
	}

	// Set the expired token cookie
	c.Cookie(&cookie)

	// Return a success response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"msg": "Logout successful",
	})
}
