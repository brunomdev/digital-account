package account

import (
	"context"
	"github.com/brunomdev/digital-account/entity"
	"github.com/pkg/errors"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, docNumber string, availableCreditLimit float64) (*entity.Account, error) {
	return s.repo.Save(ctx, docNumber, availableCreditLimit)
}

func (s *service) Get(ctx context.Context, id int) (*entity.Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateCreditLimit(ctx context.Context, id int, newLimit float64) (*entity.Account, error) {
	account, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, entity.ErrNotFound) {
		return nil, errors.Wrap(err, "account")
	}

	account.AvailabelCreditLimit = newLimit

	return s.repo.Update(ctx, account)
}
