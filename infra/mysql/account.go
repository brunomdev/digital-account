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

func (r *accountRepository) Save(ctx context.Context, docNumber string) (*entity.Account, error) {
	stmt, err := r.db.PrepareContext(ctx, `INSERT INTO accounts (document_number) VALUES(?)`)
	if err != nil {
		return nil, err
	}

	result, err := stmt.ExecContext(ctx, docNumber)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &entity.Account{
		ID:             int(id),
		DocumentNumber: docNumber,
	}, nil
}

func (r accountRepository) GetByID(ctx context.Context, id int) (*entity.Account, error) {
	stmt, err := r.db.PrepareContext(ctx, `SELECT id, document_number FROM accounts WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	var acc entity.Account
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&acc.ID, &acc.DocumentNumber)
		if err != nil {
			return nil, err
		}
	}

	if acc.ID < 1 {
		return nil, entity.ErrNotFound
	}

	return &acc, nil
}
