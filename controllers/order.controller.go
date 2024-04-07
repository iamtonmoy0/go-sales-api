package controllers

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	database "github.com/iamtonmoy0/go-sales-api/config"
	"github.com/iamtonmoy0/go-sales-api/models"
)

// create order
func CreateOrderController(c *fiber.Ctx) error {
	db := database.Database()
	// product struct
	type products struct {
		ProductId int `json:"productId"`
		Quantity  int `json:"quantity"`
	}
	body := struct {
		PaymentId int        `json:"paymentId`
		TotalPaid int        `json:"totalPaid"`
		Products  []products `json:"products"`
	}{}
	if err := c.BodyParser(&body); err != nil {
		log.Fatal("failed to parse the json")

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"msg": "filed to parse the json data"})
	}
	prodResponse := make([]*models.ProductResponseOrder, 0)

	var totalInvoicePrice = struct {
		tPrice int
	}{}

	productsIds := ""
	quantities := ""

	for _, v := range body.Products {
		totalPrice := 0
		productsIds = productsIds + "," + strconv.Itoa(v.ProductId)
		quantities = quantities + "," + strconv.Itoa(v.Quantity)

		prods := models.ProductOrder{}
		var discount models.Discount
		db.Table("products").Where("id=?", v.ProductId).First(&prods)
		db.Where("id = ?", prods.DiscountId).Find(&discount)
		discCount := 0

		if discount.Type == "BUY_N" {
			totalPrice = prods.Price * v.Quantity

			discCount = totalPrice - discount.Result
			totalInvoicePrice.tPrice = totalInvoicePrice.tPrice + discCount

		}

		if discount.Type == "PERCENT" {
			totalPrice = prods.Price * v.Quantity
			percentage := totalPrice * discount.Result / 100
			discCount = totalPrice - percentage
			totalInvoicePrice.tPrice = totalInvoicePrice.tPrice + discCount
		}

		prodResponse = append(prodResponse,
			&models.ProductResponseOrder{
				ProductId:        prods.ID,
				Name:             prods.Name,
				Price:            prods.Price,
				Discount:         discount,
				Qty:              v.Quantity,
				TotalNormalPrice: prods.Price,
				TotalFinalPrice:  discCount,
			},
		)

	}
	orderResp := models.Order{
		CashierID:      1,
		PaymentTypesId: body.PaymentId,
		TotalPrice:     totalInvoicePrice.tPrice,
		TotalPaid:      body.TotalPaid,
		TotalReturn:    body.TotalPaid - totalInvoicePrice.tPrice,
		ReceiptId:      "R000" + strconv.Itoa(rand.Intn(1000)),
		ProductId:      productsIds,
		Quantities:     quantities,
		UpdatedAt:      time.Now().UTC(),
		CreatedAt:      time.Now().UTC(),
	}
	db.Create(&orderResp)

	return c.Status(200).JSON(fiber.Map{

		"message": "success",
		"success": true,
		"data": map[string]interface{}{
			"order":    orderResp,
			"products": prodResponse,
		},
	})

}
