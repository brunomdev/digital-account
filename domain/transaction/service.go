package transaction

import (
	"context"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/domain/operationtype"
	"github.com/brunomdev/digital-account/entity"
	"github.com/pkg/errors"
	"math"
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
	acc, err := s.accountService.Get(ctx, accountID)
	if errors.Is(err, entity.ErrNotFound) {
		return nil, errors.Wrap(err, "acc")
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

	var newLimit float64
	if operationTypeID == 4 {
		if amount < 0 {
			return nil, entity.ErrInvalidAmount
		}
		newLimit = acc.AvailabelCreditLimit + amount
	} else {
		newLimit = acc.AvailabelCreditLimit - math.Abs(amount)
	}

	if newLimit <= 0 {
		return nil, entity.ErrInsufficientCreditLimit
	}

	_, err = s.accountService.UpdateCreditLimit(ctx, acc.ID, newLimit)
	if err != nil {
		return nil, errors.Wrap(err, "Create")
	}

	transaction, err := s.repo.Save(ctx, accountID, operationTypeID, amount)
	if err != nil {
		_, errUpd := s.accountService.UpdateCreditLimit(ctx, acc.ID, acc.AvailabelCreditLimit)
		if errUpd != nil {
			return nil, errors.Wrap(errUpd, "Create")
		}

		return nil, errors.Wrap(err, "Create")
	}

	return transaction, nil
}
