package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brunomdev/digital-account/entity"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_operationTypeRepository_GetByID(t *testing.T) {
	selectQuery := "SELECT id, description FROM operation_types WHERE id = ?"

	type args struct {
		ctx context.Context
		id  int
	}
	testCases := []struct {
		name    string
		mock    func() (*sql.DB, sqlmock.Sqlmock, error)
		args    args
		want    *entity.OperationType
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error prepare",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(selectQuery).
					WillReturnError(errors.New("error"))

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Error query",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(selectQuery).
					ExpectQuery().WithArgs(1).WillReturnError(errors.New("error"))

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Error row scan",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(selectQuery).
					ExpectQuery().WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "description"}).
							AddRow(false, "PAGAMENTO A VISTA"),
					)

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Error not found",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(selectQuery).
					ExpectQuery().WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "description"}),
					)

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Success",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(selectQuery).
					ExpectQuery().WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "description"}).
							AddRow(1, "PAGAMENTO A VISTA"),
					)

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want: &entity.OperationType{
				ID:          1,
				Description: "PAGAMENTO A VISTA",
			},
			wantErr: assert.NoError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := tc.mock()
			assert.NoError(t, err)

			defer func() {
				db.Close()
				assert.NoError(t, mock.ExpectationsWereMet())
			}()

			r := NewOperationTypeRepository(db)

			got, err := r.GetByID(tc.args.ctx, tc.args.id)
			if !tc.wantErr(t, err, fmt.Sprintf("GetByID(%v, %v)", tc.args.ctx, tc.args.id)) {
				return
			}
			assert.Equalf(t, tc.want, got, "GetByID(%v, %v)", tc.args.ctx, tc.args.id)
		})
	}
}
