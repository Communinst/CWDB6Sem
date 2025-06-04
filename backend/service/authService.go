package service

import (
	"context"
	"strconv"
	"time"

	authToken "github.com/Communinst/CWDB6Sem/backend/JSONWebTokens"
	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/Communinst/CWDB6Sem/backend/repository"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	repo repository.AuthRepo
}

func NewAuthService(repo repository.AuthRepo) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (service *AuthService) GenerateAuthToken(user *entities.User, secret string, expireTime int) (string, error) {
	claims := &authToken.JWTToken{
		Email: user.Email,
		Id:    strconv.Itoa(user.UserId),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expireTime))),
			Issuer:    "CWDB6Sem",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (service *AuthService) PostUser(ctx context.Context, user *entities.User) error {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	_, err := service.repo.PostUser(c, user)
	return err
}

func (service *AuthService) GetUser(ctx context.Context, userId int) (*entities.User, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return service.repo.GetUser(c, userId)
}

func (service *AuthService) GetUserByEmail(ctx context.Context, userEmail string) (*entities.User, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return service.repo.GetUserByEmail(c, userEmail)
}
