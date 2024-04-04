package models

import "time"

type Discount struct {
	ID              int       `json:"id" gorm:"autoIncrement;primaryKey"`
	Qty             int       `json:"qty"`
	Type            string    `json:"type"`
	Result          int       `json:"result"`
	ExpiredAt       int       `json:"expiredAt"`
	ExpiredAtFormat string    `json:"expiredAtFormat"`
	StringFormat    string    `json:"stringFormat"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
