package mysql

import (
	"context"
	"database/sql"
	"github.com/brunomdev/digital-account/domain/transaction"
	"github.com/brunomdev/digital-account/entity"
	"time"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) transaction.Repository {
	return &transactionRepository{db: db}
}

func (r transactionRepository) Save(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
	stmt, err := r.db.PrepareContext(ctx, `insert into transactions (account_id, operation_type_id, amount) values(?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, accountID, operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entity.Transaction{
		ID:              int(id),
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
		EventDate:       time.Now(),
	}, nil
}
