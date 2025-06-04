package repository

import (
	"context"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/jmoiron/sqlx"
)

type dumpRepo struct {
	db *sqlx.DB
}

func NewDumpRepo(db *sqlx.DB) *dumpRepo {
	return &dumpRepo{db: db}
}

func (r *dumpRepo) InsertDump(ctx context.Context, dump *entities.Dump) error {
	query := "INSERT INTO dumps (filename, size) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, dump.Filename, dump.Size)
	return err
}

func (r *dumpRepo) GetAllDumps(ctx context.Context) ([]entities.Dump, error) {
	var dumps []entities.Dump
	query := "SELECT * FROM dumps"
	err := r.db.SelectContext(ctx, &dumps, query)
	return dumps, err
}
