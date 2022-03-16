package mysql

import (
	"context"
	"database/sql"
	"github.com/brunomdev/digital-account/domain/operationtype"
	"github.com/brunomdev/digital-account/entity"
)

type operationTypeRepository struct {
	db *sql.DB
}

func NewOperationTypeRepository(db *sql.DB) operationtype.Repository {
	return &operationTypeRepository{db: db}
}

func (r operationTypeRepository) GetByID(ctx context.Context, id int) (*entity.OperationType, error) {
	stmt, err := r.db.PrepareContext(ctx, `select id, description from operation_types where id = ?`)
	if err != nil {
		return nil, err
	}

	var opType entity.OperationType
	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&opType.ID, &opType.Description)
		if err != nil {
			return nil, err
		}
	}

	if opType.ID < 1 {
		return nil, entity.ErrNotFound
	}

	return &opType, nil
}
