package transaction

import (
	"context"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/domain/account/mock_account"
	"github.com/brunomdev/digital-account/domain/operationtype"
	"github.com/brunomdev/digital-account/domain/operationtype/mock_operationtype"
	"github.com/brunomdev/digital-account/domain/transaction/mock_transaction"
	"github.com/brunomdev/digital-account/entity"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func Test_service_Create(t *testing.T) {
	type args struct {
		ctx                        context.Context
		accountID, operationTypeID int
		amount                     float64
	}
	testCases := []struct {
		name    string
		svcArgs func(ctrl *gomock.Controller) (Repository, account.Service, operationtype.Service)
		args    args
		want    *entity.Transaction
		wantErr bool
	}{
		{
			name: "Error account not found",
			svcArgs: func(ctrl *gomock.Controller) (Repository, account.Service, operationtype.Service) {
				repo := mock_transaction.NewMockRepository(ctrl)
				accountSvc := mock_account.NewMockService(ctrl)
				opTypeSvc := mock_operationtype.NewMockService(ctrl)

				accountSvc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, entity.ErrNotFound)

				return repo, accountSvc, opTypeSvc
			},
			args: args{
				accountID:       1,
				operationTypeID: 4,
				amount:          50.00,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error account service",
			svcArgs: func(ctrl *gomock.Controller) (Repository, account.Service, operationtype.Service) {
				repo := mock_transaction.NewMockRepository(ctrl)
				accountSvc := mock_account.NewMockService(ctrl)
				opTypeSvc := mock_operationtype.NewMockService(ctrl)

				accountSvc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

				return repo, accountSvc, opTypeSvc
			},
			args: args{
				accountID:       1,
				operationTypeID: 4,
				amount:          50.00,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error operation type not found",
			svcArgs: func(ctrl *gomock.Controller) (Repository, account.Service, operationtype.Service) {
				repo := mock_transaction.NewMockRepository(ctrl)
				accountSvc := mock_account.NewMockService(ctrl)
				opTypeSvc := mock_operationtype.NewMockService(ctrl)

				accountSvc.EXPECT().Get(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, accountID int) (*entity.Account, error) {
						return &entity.Account{
							ID:             accountID,
							DocumentNumber: "12345678900",
						}, nil
					})

				opTypeSvc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, entity.ErrNotFound)

				return repo, accountSvc, opTypeSvc
			},
			args: args{
				accountID:       1,
				operationTypeID: 4,
				amount:          50.00,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Error operation type service",
			svcArgs: func(ctrl *gomock.Controller) (Repository, account.Service, operationtype.Service) {
				repo := mock_transaction.NewMockRepository(ctrl)
				accountSvc := mock_account.NewMockService(ctrl)
				opTypeSvc := mock_operationtype.NewMockService(ctrl)

				accountSvc.EXPECT().Get(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, accountID int) (*entity.Account, error) {
						return &entity.Account{
							ID:             accountID,
							DocumentNumber: "12345678900",
						}, nil
					})

				opTypeSvc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repo, accountSvc, opTypeSvc
			},
			args: args{
				accountID:       1,
				operationTypeID: 4,
				amount:          50.00,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success",
			svcArgs: func(ctrl *gomock.Controller) (Repository, account.Service, operationtype.Service) {
				repo := mock_transaction.NewMockRepository(ctrl)
				accountSvc := mock_account.NewMockService(ctrl)
				opTypeSvc := mock_operationtype.NewMockService(ctrl)

				accountSvc.EXPECT().Get(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, accountID int) (*entity.Account, error) {
						return &entity.Account{
							ID:             accountID,
							DocumentNumber: "12345678900",
						}, nil
					})

				opTypeSvc.EXPECT().Get(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, operationTypeID int) (*entity.OperationType, error) {
						return &entity.OperationType{
							ID:          operationTypeID,
							Description: "PAGAMENTO",
						}, nil
					})

				repo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
						return &entity.Transaction{
							ID:              1,
							AccountID:       accountID,
							OperationTypeID: operationTypeID,
							Amount:          amount,
							EventDate:       time.Time{},
						}, nil
					})

				return repo, accountSvc, opTypeSvc
			},
			args: args{
				accountID:       1,
				operationTypeID: 4,
				amount:          50.00,
			},
			want: &entity.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          50.00,
				EventDate:       time.Time{},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := NewService(tc.svcArgs(ctrl))

			got, err := s.Create(tc.args.ctx, tc.args.accountID, tc.args.operationTypeID, tc.args.amount)
			if (err != nil) != tc.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !cmp.Equal(got, tc.want) {
				t.Errorf("Create() got = %v, want %v, %v", got, tc.want, cmp.Diff(got, tc.want))
			}
		})
	}
}
