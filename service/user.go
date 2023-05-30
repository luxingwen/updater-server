package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"updater-server/model"

	"updater-server/pkg/app"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) CreateUser(ctx *app.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := ctx.DB.Create(user).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *UserService) GetUserByID(ctx *app.Context, id uint) (*model.User, error) {
	user := &model.User{}
	err := ctx.DB.First(user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil to indicate record not found
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx *app.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	err := ctx.DB.Save(user).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *UserService) DeleteUser(ctx *app.Context, user *model.User) error {
	err := ctx.DB.Delete(user).Error
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) FindUserByUsername(ctx *app.Context, username, password string) (*model.User, error) {
	user := &model.User{}
	err := ctx.DB.Where("username = ? AND password = ?", username, password).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return user, nil
}

func (s *UserService) QueryUsers(ctx *app.Context, query *model.UserQuery) ([]model.User, error) {
	var users []model.User

	db := ctx.DB
	if query.Username != "" {
		db = db.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Email != "" {
		db = db.Where("email LIKE ?", "%"+query.Email+"%")
	}
	if query.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+query.Phone+"%")
	}

	err := db.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	return users, nil
}
