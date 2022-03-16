package transaction

import (
	"context"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/domain/operationtype"
	"github.com/brunomdev/digital-account/entity"
	"github.com/pkg/errors"
)

type service struct {
	repo           Repository
	accountService account.Service
	opTypeService  operationtype.Service
}

func NewService(repo Repository, accountService account.Service, operationTypeService operationtype.Service) Service {
	return &service{
		repo:           repo,
		accountService: accountService,
		opTypeService:  operationTypeService,
	}
}

func (s *service) Create(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
	_, err := s.accountService.Get(ctx, accountID)
	if errors.Is(err, entity.ErrNotFound) {
		return nil, errors.Wrap(err, "account")
	}
	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}

	_, err = s.opTypeService.Get(ctx, operationTypeID)
	if errors.Is(err, entity.ErrNotFound) {
		return nil, errors.Wrap(err, "operation type")
	}

	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}

	return s.repo.Save(ctx, accountID, operationTypeID, amount)
}
