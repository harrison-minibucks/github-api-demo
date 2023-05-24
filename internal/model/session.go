package model

import "time"

type Session struct {
	Id     string
	GhId   string
	Expiry time.Time
}
