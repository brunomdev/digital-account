//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=contract.go -destination=mock_account/contract.go

package account

import (
	"context"
	"github.com/brunomdev/digital-account/entity"
)

type Service interface {
	Create(ctx context.Context, docNumber string) (*entity.Account, error)
	Get(ctx context.Context, id int) (*entity.Account, error)
}

type Repository interface {
	Save(ctx context.Context, docNumber string) (*entity.Account, error)
	GetByID(ctx context.Context, id int) (*entity.Account, error)
}
