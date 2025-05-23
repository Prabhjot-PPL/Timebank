package dto

import "time"

type Session struct {
	SessionDate   time.Time `json:"session_date"`
	Helper        int       `json:"helper"`
	Recipient     int       `json:"recipient"`
	Skill         string    `json:"skill"`
	Hours         float64   `json:"hours"`
	SessionStatus string    `json:"session_status"`
}
