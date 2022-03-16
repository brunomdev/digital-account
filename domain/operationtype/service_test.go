package operationtype

import (
	"context"
	"github.com/brunomdev/digital-account/domain/operationtype/mock_operationtype"
	"github.com/brunomdev/digital-account/entity"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"testing"
)

func Test_service_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
	}
	testCases := []struct {
		name    string
		svcArgs func(ctrl *gomock.Controller) Repository
		args    args
		want    *entity.OperationType
		wantErr bool
	}{
		{
			name: "Error databasee",
			svcArgs: func(ctrl *gomock.Controller) Repository {
				repo := mock_operationtype.NewMockRepository(ctrl)

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
				repo := mock_operationtype.NewMockRepository(ctrl)

				repo.EXPECT().GetByID(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int) (*entity.OperationType, error) {
						return &entity.OperationType{
							ID:          id,
							Description: "COMPRA A VISTA",
						}, nil
					})

				return repo
			},
			args: args{
				id: 1,
			},
			want: &entity.OperationType{
				ID:          1,
				Description: "COMPRA A VISTA",
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
