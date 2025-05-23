package ports

import (
	"context"
	"timebank/src/internal/core/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user dto.UserDetails) error
	GetUserByEmail(ctx context.Context, email string) (dto.UserDetails, error)
	CreateSession(ctx context.Context, session dto.Session) error
	CompleteSessionTx(ctx context.Context, sessionID int, feedback, status string) error
}
