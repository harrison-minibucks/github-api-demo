package model

import "time"

type Item struct {
	Id          string `gorm:"primaryKey"`
	Title       string
	Description string
	Marked      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
