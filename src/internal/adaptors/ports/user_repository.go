package ports

import (
	"context"
	"timebank/src/internal/core/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user dto.UserDetails) error
	GetUserByEmail(ctx context.Context, email string) (dto.UserDetails, error)
	FindHelperBySkill(ctx context.Context, skill string) ([]dto.HelperDetails, error)
	CreateSession(ctx context.Context, session dto.Session) error
	StartSession(ctx context.Context, session_id int) error
	CompleteSessionTx(ctx context.Context, sessionID int, feedback, status string) error
}
