//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=contract.go -destination=mock_transaction/contract.go

package transaction

import (
	"context"
	"github.com/brunomdev/digital-account/entity"
)

type Service interface {
	Create(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error)
}

type Repository interface {
	Save(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error)
}
