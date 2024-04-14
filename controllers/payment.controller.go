package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
)

// Payment struct with two values
type Payment struct {
	Id            uint   `json:"paymentId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	PaymentTypeId int    `json:"payment_type_id"`
	Logo          string `json:"logo"`
}

func CreatePayment(c *fiber.Ctx) error {
	db := database.Database()
	var data map[string]string
	paymentError := c.BodyParser(&data)
	if paymentError != nil {
		log.Fatalf("payment error in post request %v", paymentError)
	}
	if data["name"] == "" || data["type"] == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment Name and Type is required",
			"error":   map[string]interface{}{},
		})
	}

	var paymentTypes models.PaymentType
	db.Where("name", data["type"]).First(&paymentTypes)
	payment := models.Payment{
		Name:          data["name"],
		Type:          data["type"],
		PaymentTypeId: int(paymentTypes.ID),
		Logo:          data["logo"],
	}
	db.Create(&payment)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    payment,
	})
}
func PaymentList(c *fiber.Ctx) error {
	db := database.Database()
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var payment []Payment
	db.Select("id ,name,type,payment_type_id,logo,created_at,updated_at").Limit(limit).Offset(skip).Find(&payment).Count(&count)
	metaMap := map[string]interface{}{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}
	categoriesData := map[string]interface{}{
		"payments": payment,
		"meta":     metaMap,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    categoriesData,
	})

}

func GetPaymentDetails(c *fiber.Ctx) error {
	db := database.Database()
	paymentId := c.Params("paymentId")

	var payment models.Payment
	db.Where("id=?", paymentId).First(&payment)

	if payment.Name == "" {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
			"error":   map[string]interface{}{},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    payment,
	})
}

func DeletePayment(c *fiber.Ctx) error {
	db := database.Database()
	paymentId := c.Params("paymentId")
	var payment models.Payment

	db.First(&payment, paymentId)
	if payment.Name == "" {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"Message": "No payment found against this payment id",
		})
	}

	result := db.Delete(&payment)
	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "payment removing failed",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
	})
}

func UpdatePayment(c *fiber.Ctx) error {
	db := database.Database()
	paymentId := c.Params("paymentId")
	fmt.Println("-----------------------------------")
	fmt.Println("---------------Params payment id--------------------", paymentId)
	fmt.Println("-----------------------------------")
	var totalPayment models.Payment
	db.Find(&totalPayment)

	fmt.Println("-----------------------------------")
	fmt.Println("---------------All payments--------------------", totalPayment)
	fmt.Println("-----------------------------------")
	var payment models.Payment

	db.Find(&payment, "id = ?", paymentId)

	var updatePaymentData models.Payment
	c.BodyParser(&updatePaymentData)
	if updatePaymentData.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Payment name is required",
			"error":   map[string]interface{}{},
		})
	}

	var paymentTypeId int
	if updatePaymentData.Type == "CASH" {
		paymentTypeId = 1
	}
	if updatePaymentData.Type == "E-WALLET" {
		paymentTypeId = 2
	}
	if updatePaymentData.Type == "EDC" {
		paymentTypeId = 3
	}
	fmt.Println(paymentTypeId)
	payment.Name = updatePaymentData.Name
	payment.Type = updatePaymentData.Type
	payment.PaymentTypeId = paymentTypeId
	payment.Logo = updatePaymentData.Logo

	db.Save(&payment)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"Message": "Success",
		"data":    payment,
	})

}
