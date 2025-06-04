package service

import (
	"context"

	entities "github.com/Communinst/CWDB6Sem/backend/entity"
	"github.com/Communinst/CWDB6Sem/backend/repository"
)

type DumpService struct {
	repo repository.DumpRepo
}

func NewDumpService(repo repository.DumpRepo) *DumpService {
	return &DumpService{repo: repo}
}

func (s *DumpService) InsertDump(ctx context.Context, filePath string, size int64) error {
	dump := &entities.Dump{
		Filename: filePath,
		Size:     size,
	}
	return s.repo.InsertDump(ctx, dump)
}

func (s *DumpService) GetAllDumps(ctx context.Context) ([]entities.Dump, error) {
	return s.repo.GetAllDumps(ctx)
}
