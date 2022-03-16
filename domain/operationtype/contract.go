//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=contract.go -destination=mock_operationtype/contract.go

package operationtype

import (
	"context"
	"github.com/brunomdev/digital-account/entity"
)

type Service interface {
	Get(ctx context.Context, id int) (*entity.OperationType, error)
}

type Repository interface {
	GetByID(ctx context.Context, id int) (*entity.OperationType, error)
}
