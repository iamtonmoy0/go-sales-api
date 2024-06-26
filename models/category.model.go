package models

import "time"

type Category struct {
	ID       int       `json:"id" gorm:"autoIncrement;primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
