package storage

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	cnfg "github.com/Communinst/CWDB6Sem/backend/config"
	customErrors "github.com/Communinst/CWDB6Sem/backend/errors"
	"github.com/jmoiron/sqlx"
)

func InitDBConn(config *cnfg.Database) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.DBName,
		config.SSLMode)

	//slog.Info("DB Connection String:", connStr)  // Debug log

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to open database connection: ", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to ping database: ", err)
		return nil
	}
	log.Print("Successfully connected to the database!")
	return db
}

func RunDBTableScript(db *sqlx.DB, scriptPath string) error {

	script, err := os.ReadFile(scriptPath)
	if err != nil {
		log.Fatalf("Failed to read SQL script file: %v", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		slog.Error("transaction initiation error")
		return &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        "transaction initiation failed",
		}
	}

	_, err = db.Exec(string(script))
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to execute SQL script: %v", err)
		return &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        fmt.Sprintf("failed to execute SQL script down the path: %s", scriptPath),
		}
	}

	if err = tx.Commit(); err != nil {
		slog.Error("transaction fulfillment error")
		return &customErrors.ErrorWithStatusCode{
			HTTPStatus: http.StatusInternalServerError,
			Msg:        "transaction fulfillment failed",
		}
	}

	//log.Print("Script down the path:", scriptPath, ": succesfull run")
	return nil
}

func CloseDBConn(db *sqlx.DB) error {
	return db.Close()
}
