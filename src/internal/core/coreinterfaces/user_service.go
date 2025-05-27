package coreinterfaces

import (
	"context"
	"timebank/src/internal/core/dto"
)

type Service interface {
	RegisterUser(ctx context.Context, requestData dto.UserDetails) error
	LoginUser(ctx context.Context, requestData dto.UserDetails) (dto.LoginResponse, error)
	FindUsersBySkill(ctx context.Context, skill string) ([]dto.HelperDetails, error)
	CreateSession(ctx context.Context, session dto.Session) error
	StartSession(ctx context.Context, session_id int) error
	CompleteSession(ctx context.Context, sessionID int, feedback string, status string) error
}
