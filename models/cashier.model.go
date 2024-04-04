package models

import "time"

type Cashier struct {
	ID        int       `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Password  string    `json:"password,omitempty" gorm:"not null;size:256"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
