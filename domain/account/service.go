package account

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

func (s *service) Create(ctx context.Context, docNumber string) (*entity.Account, error) {
	return s.repo.Save(ctx, docNumber)
}

func (s *service) Get(ctx context.Context, id int) (*entity.Account, error) {
	return s.repo.GetByID(ctx, id)
}
