package controller

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"updater-server/model"
	"updater-server/pkg/app"
	"updater-server/service"
)

type AuthController struct {
	AuthService *service.AuthService
}

func (cc *AuthController) Login(ctx *app.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	user, err := cc.AuthService.Login(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSONError(http.StatusForbidden, err.Error())
		return
	}

	// Generate JWT token
	token, err := generateJWTToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSONSuccess(gin.H{"token": token})
}

func generateJWTToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"uuid":     user.Uuid,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (1 day)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-secret-key")) // Replace with your secret key
}
