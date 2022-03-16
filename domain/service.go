package domain

import (
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/domain/operationtype"
	"github.com/brunomdev/digital-account/domain/transaction"
)

type Service struct {
	Account       account.Service
	OperationType operationtype.Service
	Transaction   transaction.Service
}
