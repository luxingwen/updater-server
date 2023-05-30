package service

import (
	"fmt"

	"gorm.io/gorm"

	"updater-server/model"
	"updater-server/pkg/app"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Login(ctx *app.Context, username, password string) (*model.User, error) {
	user, err := s.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.Password != password {
		return nil, fmt.Errorf("incorrect password")
	}

	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	if err == bcrypt.ErrMismatchedHashAndPassword {
	// 		return nil, fmt.Errorf("incorrect password")
	// 	}
	// 	return nil, fmt.Errorf("failed to compare password: %w", err)
	// }

	return user, nil
}

func (s *AuthService) FindUserByUsername(ctx *app.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := ctx.DB.Where("username = ?", username).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}
