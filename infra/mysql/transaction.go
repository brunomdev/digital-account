package mysql

import (
	"context"
	"database/sql"
	"github.com/brunomdev/digital-account/domain/transaction"
	"github.com/brunomdev/digital-account/entity"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) transaction.Repository {
	return &transactionRepository{db: db}
}

func (r transactionRepository) Save(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO transactions (account_id, operation_type_id, amount) VALUES(?, ?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, accountID, operationTypeID, amount)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, int(id))
}

func (r transactionRepository) GetByID(ctx context.Context, id int) (*entity.Transaction, error) {
	stmt, err := r.db.PrepareContext(ctx, `SELECT id, account_id, operation_type_id, amount, created_at FROM transactions WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	var txn entity.Transaction
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&txn.ID, &txn.AccountID, &txn.OperationTypeID, &txn.Amount, &txn.EventDate)
		if err != nil {
			return nil, err
		}
	}

	if txn.ID < 1 {
		return nil, entity.ErrNotFound
	}

	return &txn, nil
}
