package account

import (
	"context"
	"github.com/brunomdev/digital-account/domain/account/mock_account"
	"github.com/brunomdev/digital-account/entity"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"testing"
)

func Test_service_Create(t *testing.T) {
	type args struct {
		ctx                  context.Context
		docNumber            string
		availableCreditLimit float64
	}
	testCases := []struct {
		name    string
		svcArgs func(ctrl *gomock.Controller) Repository
		args    args
		want    *entity.Account
		wantErr bool
	}{
		{
			name: "Error database",
			svcArgs: func(ctrl *gomock.Controller) Repository {
				repo := mock_account.NewMockRepository(ctrl)

				repo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repo
			},
			args: args{
				docNumber: "12345678900",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success",
			svcArgs: func(ctrl *gomock.Controller) Repository {
				repo := mock_account.NewMockRepository(ctrl)

				repo.EXPECT().Save(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, docNumber string, availableCreditLimit float64) (*entity.Account, error) {
						return &entity.Account{
							ID:                   1,
							DocumentNumber:       docNumber,
							AvailabelCreditLimit: availableCreditLimit,
						}, nil
					})

				return repo
			},
			args: args{
				docNumber: "12345678900",
			},
			want: &entity.Account{
				ID:             1,
				DocumentNumber: "12345678900",
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := NewService(tc.svcArgs(ctrl))

			got, err := s.Create(tc.args.ctx, tc.args.docNumber, tc.args.availableCreditLimit)
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

func Test_service_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	testCases := []struct {
		name    string
		svcArgs func(ctrl *gomock.Controller) Repository
		args    args
		want    *entity.Account
		wantErr bool
	}{
		{
			name: "Error database",
			svcArgs: func(ctrl *gomock.Controller) Repository {
				repo := mock_account.NewMockRepository(ctrl)
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

				return repo
			},
			args: args{
				id: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Success",
			svcArgs: func(ctrl *gomock.Controller) Repository {
				repo := mock_account.NewMockRepository(ctrl)
				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int) (*entity.Account, error) {
						return &entity.Account{
							ID:                   id,
							DocumentNumber:       "12345678900",
							AvailabelCreditLimit: 5000.00,
						}, nil
					})

				return repo
			},
			args: args{
				id: 1,
			},
			want: &entity.Account{
				ID:                   1,
				DocumentNumber:       "12345678900",
				AvailabelCreditLimit: 5000.00,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := NewService(tc.svcArgs(ctrl))

			got, err := s.Get(tc.args.ctx, tc.args.id)
			if (err != nil) != tc.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !cmp.Equal(got, tc.want) {
				t.Errorf("Get() got = %v, want %v, %v", got, tc.want, cmp.Diff(got, tc.want))
			}
		})
	}
}
