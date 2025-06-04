package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	customErrors "github.com/Communinst/CWDB6Sem/backend/errors"
	"github.com/jmoiron/sqlx"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *authRepo {
	return &authRepo{
		db: db,
	}
}

func (repo *authRepo) PostUser(ctx context.Context, user *entities.User) (int, error) {
	var result_id int

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("transaction initiation error")
		return -1, &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        "transaction initiation failed",
		}
	}

	query := fmt.Sprintf(`INSERT INTO %s (login, password, nickname, email, sign_up_date, role_id)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id`, usersTable)

	err = tx.QueryRowContext(ctx, query,
		user.Login,
		user.Password,
		user.Nickname,
		user.Email,
		user.SignUpDate,
		user.RoleId).Scan(&result_id)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		slog.Error("transaction fulfillment error")
		return -1, &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        "transaction fulfillment failed",
		}
	}

	log.Printf("User with email %s posted successfully", user.Email)
	return result_id, err
}

func (repo *authRepo) GetUser(ctx context.Context, userId int) (*entities.User, error) {
	var resultUser entities.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE user_id = $1`, usersTable)

	err := repo.db.GetContext(ctx, &resultUser, query, userId)
	if err == nil {
		log.Printf("User by %d id was obtained", userId)
		return &resultUser, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &resultUser, &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("user with %d if wasn't found", userId),
		}
	}

	slog.Error("unknown error obtaining user by id")
	return &resultUser, &customErrors.ErrorWithStatusCode{
		HTTPStatus: http.StatusInternalServerError,
		Msg:        "unknown interanal server error occured",
	}
}

func (repo *authRepo) GetUserByEmail(ctx context.Context, userEmail string) (*entities.User, error) {
	var resultUser entities.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE email = $1`, usersTable)

	fmt.Printf("%s", userEmail)
	err := repo.db.GetContext(ctx, &resultUser, query, userEmail)
	if err == nil {
		log.Printf("User by %s email was obtained", userEmail)
		return &resultUser, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return &resultUser, &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("user with %s email wasn't found", userEmail),
		}
	}

	slog.Error("unknown error obtaining user by id")
	return &resultUser, &customErrors.ErrorWithStatusCode{
		HTTPStatus: http.StatusInternalServerError,
		Msg:        "unknown interanal server error occured",
	}
}
