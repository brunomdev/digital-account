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
)

func Test_accountRepository_Save(t *testing.T) {
	insertQuery := "INSERT INTO accounts (document_number, available_credit_limit) VALUES(?, ?)"

	type args struct {
		ctx                  context.Context
		docNumber            string
		availableCreditLimit float64
	}
	testCases := []struct {
		name    string
		mock    func() (*sql.DB, sqlmock.Sqlmock, error)
		args    args
		want    *entity.Account
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
				ctx:                  context.TODO(),
				docNumber:            "12345678900",
				availableCreditLimit: 50.00,
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
					WithArgs("12345678900", 50.00).
					WillReturnError(errors.New("error"))

				return db, mock, nil
			},
			args: args{
				ctx:                  context.TODO(),
				docNumber:            "12345678900",
				availableCreditLimit: 50.00,
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
				ctx:       context.TODO(),
				docNumber: "12345678900",
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
					WillReturnResult(sqlmock.NewResult(1, 1))

				return db, mock, nil
			},
			args: args{
				ctx:                  context.TODO(),
				docNumber:            "12345678900",
				availableCreditLimit: 50.00,
			},
			want: &entity.Account{
				ID:                   1,
				DocumentNumber:       "12345678900",
				AvailabelCreditLimit: 50.00,
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

			r := NewAccountRepository(db)

			got, err := r.Save(tc.args.ctx, tc.args.docNumber, tc.args.availableCreditLimit)
			if !tc.wantErr(t, err, fmt.Sprintf("Save(%v, %v)", tc.args.ctx, tc.args.docNumber)) {
				return
			}
			assert.Equalf(t, tc.want, got, "Save(%v, %v)", tc.args.ctx, tc.args.docNumber)
		})
	}
}

func Test_accountRepository_GetByID(t *testing.T) {
	selectQuery := "SELECT id, document_number, available_credit_limit FROM accounts WHERE id = ?"

	type args struct {
		ctx context.Context
		id  int
	}
	testCases := []struct {
		name    string
		mock    func() (*sql.DB, sqlmock.Sqlmock, error)
		args    args
		want    *entity.Account
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
						sqlmock.NewRows([]string{"id", "document_number", "available_credit_limit"}).
							AddRow(false, "12345678900", 50.00),
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
						sqlmock.NewRows([]string{"id", "document_number", "available_credit_limit"}),
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
						sqlmock.NewRows([]string{"id", "document_number", "available_credit_limit"}).
							AddRow(1, "12345678900", 50.00),
					)

				return db, mock, nil
			},
			args: args{
				ctx: context.TODO(),
				id:  1,
			},
			want: &entity.Account{
				ID:                   1,
				DocumentNumber:       "12345678900",
				AvailabelCreditLimit: 50.00,
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

			r := NewAccountRepository(db)

			got, err := r.GetByID(tc.args.ctx, tc.args.id)
			if !tc.wantErr(t, err, fmt.Sprintf("GetByID(%v, %v)", tc.args.ctx, tc.args.id)) {
				return
			}
			assert.Equalf(t, tc.want, got, "GetByID(%v, %v)", tc.args.ctx, tc.args.id)
		})
	}
}
