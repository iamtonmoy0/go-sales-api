package controllers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

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

// check order
func CheckOrder(c *fiber.Ctx) error {
	db := database.Database()
	param := c.Params("orderId")

	var order models.Order
	db.Where("id=?", param).First(&order)
	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  false,
			"message": "Order does not exist",
		})
	}

	if order.IsDownload == 0 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "success",
			"data": map[string]interface{}{
				"isDownloaded": false,
			},
		})
	}

	if order.IsDownload == 1 {
		return c.Status(200).JSON(fiber.Map{
			"status":  true,
			"message": "success",
			"data": map[string]interface{}{
				"isDownloaded": true,
			},
		})
	}

	return nil

}

// order details
func OrderDetail(c *fiber.Ctx) error {
	db := database.Database()

	param := c.Params("orderId")

	var order models.Order
	db.Where("id=?", param).First(&order)

	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Not Found",
			"error":   map[string]interface{}{},
		})
	}
	productIds := strings.Split(order.ProductId, ",")
	TotalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := models.Product{}
		db.Where("id = ?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := models.Cashier{}
	db.Where("id = ?", order.CashierID).Find(&cashier)

	paymentType := models.PaymentType{}
	db.Where("id = ?", order.PaymentTypesId).Find(&paymentType)

	orderTable := models.Order{}
	db.Where("id = ?", order.ID).Find(&orderTable)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": map[string]interface{}{
			"order": map[string]interface{}{
				"orderId":        order.ID,
				"cashiersId":     order.CashierID,
				"paymentTypesId": order.PaymentTypesId,
				"totalPrice":     order.TotalPrice,
				"totalPaid":      order.TotalPaid,
				"totalReturn":    order.TotalReturn,
				"receiptId":      order.ReceiptId,
				"createdAt":      order.CreatedAt,
				"cashier":        cashier,
				"payment_type":   paymentType,
			},
			"products": TotalProducts,
		},
		"Message": "Success",
	})

}

func OrdersList(c *fiber.Ctx) error {
	db := database.Database()
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var order []models.Order

	db.Select("*").Limit(limit).Offset(skip).Find(&order).Count(&count)

	type OrderList struct {
		OrderId        int                `json:"orderId"`
		CashierID      int                `json:"cashiersId"`
		PaymentTypesId int                `json:"paymentTypesId"`
		TotalPrice     int                `json:"totalPrice"`
		TotalPaid      int                `json:"totalPaid"`
		TotalReturn    int                `json:"totalReturn"`
		ReceiptId      string             `json:"receiptId"`
		CreatedAt      time.Time          `json:"createdAt"`
		Payments       models.PaymentType `json:"payment_type"`
		Cashiers       models.Cashier     `json:"cashier"`
	}
	OrderResponse := make([]*OrderList, 0)

	for _, v := range order {
		cashier := models.Cashier{}
		db.Where("id = ?", v.CashierID).Find(&cashier)
		paymentType := models.PaymentType{}
		db.Where("id = ?", v.PaymentTypesId).Find(&paymentType)

		OrderResponse = append(OrderResponse, &OrderList{
			OrderId:        v.ID,
			CashierID:      v.CashierID,
			PaymentTypesId: v.PaymentTypesId,
			TotalPrice:     v.TotalPrice,
			TotalPaid:      v.TotalPaid,
			TotalReturn:    v.TotalReturn,
			ReceiptId:      v.ReceiptId,
			CreatedAt:      v.CreatedAt,
			Payments:       paymentType,
			Cashiers:       cashier,
		})

	}

	return c.Status(404).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    OrderResponse,
		"meta": map[string]interface{}{
			"total": count,
			"limit": limit,
			"skip":  skip,
		},
	})
}

func DownloadOrder(c *fiber.Ctx) error {
	db := database.Database()
	param := c.Params("orderId")

	var order models.Order
	db.Where("id=?", param).First(&order)

	if order.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"Message": "Order Not Found",
			"error":   map[string]interface{}{},
		})
	}
	productIds := strings.Split(order.ProductId, ",")

	TotalProducts := make([]*models.Product, 0)

	for i := 1; i < len(productIds); i++ {
		prods := models.Product{}
		db.Where("id = ?", productIds[i]).Find(&prods)
		TotalProducts = append(TotalProducts, &prods)
	}
	cashier := models.Cashier{}
	db.Where("id = ?", order.CashierID).Find(&cashier)
	paymentType := models.PaymentType{}

	db.Where("id = ?", order.PaymentTypesId).Find(&paymentType)
	orderTable := models.Order{}
	db.Where("id = ?", order.ID).Find(&orderTable)

	///pdf Generating
	twoDarray := [][]string{{}}
	quantities := strings.Split(order.Quantities, ",")
	quantities = quantities[1:]
	for i := 0; i < len(TotalProducts); i++ {

		s_array := []string{}
		s_array = append(s_array, TotalProducts[i].Sku)

		s_array = append(s_array, TotalProducts[i].Name)
		s_array = append(s_array, quantities[i])
		s_array = append(s_array, strconv.Itoa(TotalProducts[i].Price))
		twoDarray = append(twoDarray, s_array)

	}

	begin := time.Now()
	grayColor := getGrayColor()
	whiteColor := color.NewWhite()
	header := getHeader()
	contents := twoDarray

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	//m.SetBorder(true)

	//Top Heading
	m.SetBackgroundColor(grayColor)
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Order Invoice #"+strconv.Itoa(order.ID), props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
			})
		})
	})
	m.SetBackgroundColor(whiteColor)

	//Table setting
	m.TableList(header, contents, props.TableList{
		HeaderProp: props.TableListContent{
			Size:      9,
			GridSizes: []uint{3, 4, 2, 3},
		},
		ContentProp: props.TableListContent{
			Size:      8,
			GridSizes: []uint{3, 4, 2, 3},
		},
		Align:                consts.Center,
		AlternatedBackground: &grayColor,
		HeaderContentSpace:   1,
		Line:                 false,
	})
	//Total price
	m.Row(20, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPrice), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})
	m.Row(21, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Paid:", props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalPaid), props.Text{
				Top:   0.5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	m.Row(22, func() {
		m.ColSpace(7)
		m.Col(2, func() {
			m.Text("Total Return", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Right,
			})
		})
		m.Col(3, func() {
			m.Text("RS. "+strconv.Itoa(order.TotalReturn), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  8,
				Align: consts.Center,
			})
		})
	})

	//Invoice creation
	currentTime := time.Now()
	pdfFileName := "invoice-" + currentTime.Format("2006-Jan-02")
	err := m.OutputFileAndClose(pdfFileName + ".pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))

	//update recept is downloaded to 1 means true
	db.Table("orders").Where("id=?", order.ID).Update("is_download", 1)
	return c.Status(200).JSON(fiber.Map{
		"Success": true,
		"Message": "Success",
	})

}

func getHeader() []string {
	return []string{"Product Sku", "Name", "Qty", "Price"}
}

func getGrayColor() color.Color {
	return color.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}
