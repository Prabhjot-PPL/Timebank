package dto

import (
	"time"
)

type LoginResponse struct {
	TokenString string
	TokenExpire time.Time
	Session     struct {
		Id        string
		ExpiresAt time.Time
	}
	FoundUser UserDetails
}
