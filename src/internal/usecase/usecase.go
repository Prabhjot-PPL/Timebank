package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"timebank/src/internal/adaptors/ports"
	"timebank/src/internal/core/coreinterfaces"
	"timebank/src/internal/core/dto"
	"timebank/src/pkg"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) coreinterfaces.Service {
	return &UserService{userRepo: userRepo}
}

// --------------------------------AUTH----------------------------------------

// REGISTER
func (u *UserService) RegisterUser(ctx context.Context, requestData dto.UserDetails) error {
	err := u.userRepo.CreateUser(ctx, requestData)
	return err
}

// LOGIN
type LoginResponse struct {
	TokenString string
	TokenExpire time.Time
	FoundUser   string
}

func (u *UserService) LoginUser(ctx context.Context, requestData dto.UserDetails) (dto.LoginResponse, error) {
	dbUser, err := u.userRepo.GetUserByEmail(ctx, requestData.Email)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	err = pkg.CheckPassword(dbUser.Password, requestData.Password)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unable to match password: %w", err)
	}

	// Create JWT token
	tokenExpire := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id": dbUser.Id,
		"email":   dbUser.Email,
		"exp":     tokenExpire.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_secret_key"))
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return dto.LoginResponse{
		TokenString: tokenString,
		TokenExpire: tokenExpire,
		FoundUser:   dbUser,
	}, nil

}

// --------------------------------SESSION----------------------------------------

func (u *UserService) CreateSession(ctx context.Context, session dto.Session) error {
	// You can add business validations here (optional)
	if session.Helper == session.Recipient {
		return errors.New("helper and recipient cannot be the same")
	}

	return u.userRepo.CreateSession(ctx, session)
}

func (u *UserService) CompleteSession(ctx context.Context, sessionID int, feedback string, status string) error {
	return u.userRepo.CompleteSessionTx(ctx, sessionID, feedback, status)
}
