package mysql

import (
	"context"
	"database/sql"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/entity"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) account.Repository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Save(ctx context.Context, docNumber string, availableCreditLimit float64) (*entity.Account, error) {
	stmt, err := r.db.PrepareContext(ctx, `insert into accounts (document_number, available_credit_limit) values(?, ?)`)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, docNumber, availableCreditLimit)
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

	return &entity.Account{
		ID:                   int(id),
		DocumentNumber:       docNumber,
		AvailabelCreditLimit: availableCreditLimit,
	}, nil
}

func (r accountRepository) GetByID(ctx context.Context, id int) (*entity.Account, error) {
	stmt, err := r.db.PrepareContext(ctx, `select id, document_number, available_credit_limit from accounts where id = ?`)
	if err != nil {
		return nil, err
	}

	var acc entity.Account
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&acc.ID, &acc.DocumentNumber, &acc.AvailabelCreditLimit)
		if err != nil {
			return nil, err
		}
	}

	if acc.ID < 1 {
		return nil, entity.ErrNotFound
	}

	return &acc, nil
}

func (r accountRepository) Update(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	stmt, err := r.db.PrepareContext(ctx, `update accounts set document_number = ?, available_credit_limit = ? where id = ?`)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, account.DocumentNumber, account.AvailabelCreditLimit, account.ID)
	if err != nil {
		return nil, err
	}

	err = stmt.Close()
	if err != nil {
		return nil, err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return nil, err
	}

	return account, nil
}
