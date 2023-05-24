package model

import "time"

type Session struct {
	Id     string    `json:"id"`
	GhId   uint32    `json:"gh_id"`
	Expiry time.Time // TODO: Implement expiry check
}

type SessionKey string
