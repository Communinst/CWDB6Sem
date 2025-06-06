package repository

import (
	"context"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/jmoiron/sqlx"
)

const (
	rolesTable = "roles"
	usersTable = "users"
)


type UserRepo interface {
	PostUser(ctx context.Context, user *entities.User) (int, error)
	GetUser(ctx context.Context, userId int) (*entities.User, error)
	GetAllUsers(ctx context.Context) ([]entities.User, error)
	DeleteUser(ctx context.Context, userId int) error
	PutUserRole(ctx context.Context, userId int, roleId int) error
}
type AuthRepo interface {
	PostUser(ctx context.Context, user *entities.User) (int, error)
	GetUser(ctx context.Context, userId int) (*entities.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (*entities.User, error)
}

type DumpRepo interface {
	InsertDump(ctx context.Context, dump *entities.Dump) error
	GetAllDumps(ctx context.Context) ([]entities.Dump, error)
}

type Repository struct {
	// RoleRepo
	UserRepo
	AuthRepo
	DumpRepo
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepo: NewUserRepo(db),
		AuthRepo: NewAuthRepo(db),
		DumpRepo: NewDumpRepo(db),
	}
}
