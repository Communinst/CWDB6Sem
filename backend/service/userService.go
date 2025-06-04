package service

import (
	"context"
	"time"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/Communinst/CWDB6Sem/backend/repository"
)

type UserService struct {
	repo repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (service *UserService) PostUser(ctx context.Context, user *entities.User) error {

	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	_, err := service.repo.PostUser(c, user)
	return err
}

func (service *UserService) GetUser(ctx context.Context, userId int) (*entities.User, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return service.repo.GetUser(c, userId)
}

func (service *UserService) GetAllUsers(ctx context.Context) ([]entities.User, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return service.repo.GetAllUsers(c)
}

func (service *UserService) DeleteUser(ctx context.Context, userId int) error {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return service.repo.DeleteUser(c, userId)
}

func (service *UserService) PutUserRole(ctx context.Context, userId int, roleId int) error {
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	return service.repo.PutUserRole(c, userId, roleId)
}
