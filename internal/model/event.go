package model

import "time"

type AddNewEventParams struct {
	Name        string    `json:"name" form:"name"`
	Description string    `json:"description" form:"description"`
	ImageUrl    string    `json:"image"`
	Start       time.Time `json:"start" form:"start"` // ISO 8601: "2025-04-08T17:30:00Z"
	End         time.Time `json:"end" form:"end"`     // ISO 8601
}

type UpdateEventParams struct {
	Name        string    `json:"name" form:"name"`
	Description string    `json:"description" form:"description"`
	ImageUrl    string    `json:"image"`
	Start       time.Time `json:"start" form:"start"`
	End         time.Time `json:"end" form:"end"`
	Active      bool      `json:"active" form:"active"`
	UserId      string    `json:"user_id"`
}

type EventQuery struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
