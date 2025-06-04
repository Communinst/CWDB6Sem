// Describe overall entities and functions
// Get back to HANDLER_MAIN to determine the relevant routes for each endpoints
// If not done inside 24 hours, you're better be dead
// Freaking idiot

// Up handler

package service

import (
	"context"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/Communinst/CWDB6Sem/backend/repository"
)

type Genre interface {
}

type Role interface {
}
type UserServiceInterface interface {
	PostUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, userId int) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	DeleteUser(ctx context.Context, userId int) error
	PutUserRole(ctx context.Context, userId int, roleId int) error
}
type AuthServiceInterface interface {
	GenerateAuthToken(user *entities.User, secret string, expireTime int) (string, error)
	PostUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, userId int) (*entities.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*entities.User, error)
}

type DumpServiceInterface interface {
	InsertDump(ctx context.Context, filePath string, size int64) error
	GetAllDumps(ctx context.Context) ([]entities.Dump, error)
}

type Service struct {
	UserServiceInterface
	AuthServiceInterface

	DumpServiceInterface
}

func New(repo *repository.Repository) *Service {
	return &Service{
		UserServiceInterface: NewUserService(repo.UserRepo),
		AuthServiceInterface: NewAuthService(repo.AuthRepo),
		DumpServiceInterface: NewDumpService(repo.DumpRepo),
	}
}
