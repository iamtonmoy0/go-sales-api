package models

import "time"

type Product struct {
	ID          int       `json:"id" gorm:"autoIncrement;primaryKey"`
	Sku         string    `json:"sku"`
	Name        string    `gorm:"column:name;not null;" json:"name"`
	Price       int       `json:"price"`
	Stock       int       `json:"stock"`
	Image       string    `json:"image"`
	FinalPrice  string    `json:"final_price"`  // Final price
	NormalPrice string    `json:"normal_price"` // Normal  price, used for compare with final price. If it is empty, use Price as normal price. after applying discounts and taxes. after discounts and taxes.
	CategoryId  string    `json:"category_id"`  // Category id
	DiscountId  string    `json:"discount_id"`  // Discount id
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at`
}

// product result
type ProductResult struct {
	Id       int      `json:"productId" gorm:"autoIncrement;primaryKey"`
	Sku      string   `json:"sku"`
	Name     string   `json:"name"`
	Stock    int      `json:"stock"`
	Price    int      `json:"price"`
	Image    string   `json:"image"`
	Category Category `json:"category"`
	Discount Discount `json:"discount"`
}
