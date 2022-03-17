package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brunomdev/digital-account/entity"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_transactionRepository_Save(t *testing.T) {
	insertQuery := "INSERT INT transactions (account_id, operation_type_id, amount) VALUES(?, ?, ?)"
	selectQuery := "SELECT id, account_id, operation_type_id, amount, created_at FROM accounts WHERE id = ?"

	type args struct {
		ctx             context.Context
		accountID       int
		operationTypeID int
		amount          float64
	}
	testCases := []struct {
		name    string
		mock    func() (*sql.DB, sqlmock.Sqlmock, error)
		args    args
		want    *entity.Transaction
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Error prepare",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(insertQuery).
					WillReturnError(errors.New("error"))

				return db, mock, nil
			},
			args: args{
				ctx:             context.TODO(),
				accountID:       1,
				operationTypeID: 1,
				amount:          123.45,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Error execution",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(insertQuery).ExpectExec().
					WithArgs(1, 1, 123.45).
					WillReturnError(errors.New("error"))

				return db, mock, nil
			},
			args: args{
				ctx:             context.TODO(),
				accountID:       1,
				operationTypeID: 1,
				amount:          123.45,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "Error without LastInsertId",
			mock: func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
				if err != nil {
					return nil, nil, err
				}

				mock.ExpectPrepare(insertQuery).ExpectExec().
					WillReturnResult(sqlmock.NewErrorResult(errors.New("error")))

				return db, mock, nil
			},
			args: args{
				ctx:             context.TODO(),
				accountID:       1,
				operationTypeID: 1,
				amount:          123.45,
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

				mock.ExpectPrepare(insertQuery).ExpectExec().
					WithArgs(1, 4, 123.45).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectPrepare(selectQuery).ExpectQuery().
					WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "account_id", "operation_type_id", "amount", "created_at"}).
							AddRow(1, 1, 4, 123.45, time.Date(2022, 3, 17, 17, 30, 0, 0, time.UTC)),
					)

				return db, mock, nil
			},
			args: args{
				ctx:             context.TODO(),
				accountID:       1,
				operationTypeID: 4,
				amount:          123.45,
			},
			want: &entity.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          123.45,
				EventDate:       time.Date(2022, 3, 17, 17, 30, 0, 0, time.UTC),
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

			r := NewTransactionRepository(db)

			got, err := r.Save(tc.args.ctx, tc.args.accountID, tc.args.operationTypeID, tc.args.amount)
			if !tc.wantErr(t, err, fmt.Sprintf(
				"Save(%v, %v, %v, %v)", tc.args.ctx, tc.args.accountID, tc.args.operationTypeID, tc.args.amount,
			)) {
				return
			}
			assert.Equalf(
				t,
				tc.want,
				got,
				"Save(%v, %v, %v, %v)",
				tc.args.ctx, tc.args.accountID, tc.args.operationTypeID, tc.args.amount,
			)
		})
	}
}

func Test_transactionRepository_GetByID(t *testing.T) {
	selectQuery := "SELECT id, account_id, operation_type_id, amount, created_at FROM accounts WHERE id = ?"

	type args struct {
		ctx context.Context
		id  int
	}
	testCases := []struct {
		name    string
		mock    func() (*sql.DB, sqlmock.Sqlmock, error)
		args    args
		want    *entity.Transaction
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
						sqlmock.NewRows([]string{"id", "acount_id", "operation_type_id", "amount", "created_at"}).
							AddRow(1, 1, 4, 123.45, "2022"),
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
						sqlmock.NewRows([]string{"id", "acount_id", "operation_type_id", "amount", "created_at"}),
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
						sqlmock.NewRows([]string{"id", "acount_id", "operation_type_id", "amount", "created_at"}).
							AddRow(1, 1, 4, 123.45, time.Date(2022, 3, 17, 17, 30, 0, 0, time.UTC)),
					)

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want: &entity.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          123.45,
				EventDate:       time.Date(2022, 3, 17, 17, 30, 0, 0, time.UTC),
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

			r := NewTransactionRepository(db)

			got, err := r.GetByID(tc.args.ctx, tc.args.id)
			if !tc.wantErr(t, err, fmt.Sprintf("GetByID(%v, %v)", tc.args.ctx, tc.args.id)) {
				return
			}
			assert.Equalf(t, tc.want, got, "GetByID(%v, %v)", tc.args.ctx, tc.args.id)
		})
	}
}
