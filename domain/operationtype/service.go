package operationtype

import (
	"context"
	"github.com/brunomdev/digital-account/entity"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Get(ctx context.Context, id int) (*entity.OperationType, error) {
	return s.repo.GetByID(ctx, id)
}
