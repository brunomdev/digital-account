package handlers

import (
	"context"
	"encoding/json"
	"github.com/brunomdev/digital-account/app/api/presenter"
	"github.com/brunomdev/digital-account/domain/transaction"
	"github.com/brunomdev/digital-account/domain/transaction/mock_transaction"
	"github.com/brunomdev/digital-account/entity"
	testHelper "github.com/brunomdev/digital-account/pkg/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func Test_transactionHandler_Create(t *testing.T) {
	testCases := []struct {
		name       string
		svcArgs    func(ctrl *gomock.Controller) transaction.Service
		reqBody    []byte
		wantStatus int
		wantBody   func() ([]byte, error)
	}{
		{
			name: "Error bodyParser",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				return mock_transaction.NewMockService(ctrl)
			},
			reqBody:    []byte(`{"acount_id": 1, "operation_type_id": 1, "amount": 123.45,}`),
			wantStatus: http.StatusBadRequest,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{
					Title:  "Unable to parse body",
					Detail: "invalid character '}' looking for beginning of value",
				})
			},
		},
		{
			name: "Error validation",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				return mock_transaction.NewMockService(ctrl)
			},
			reqBody:    []byte(`{"account_id": 1}`),
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: func() ([]byte, error) {
				return json.Marshal([]presenter.ErrorResponse{
					{
						Source: "OperationTypeID",
						Detail: "OperationTypeID is a required field",
					},
					{
						Source: "Amount",
						Detail: "Amount is a required field",
					},
				})
			},
		},
		{
			name: "Error service not found",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				svc := mock_transaction.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.Wrap(entity.ErrNotFound, "account"))

				return svc
			},
			reqBody:    []byte(`{"account_id": 1, "operation_type_id": 1, "amount": 123.45}`),
			wantStatus: http.StatusNotFound,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{
					Title:  "Resource not found",
					Detail: "account: not found",
				})
			},
		},
		{
			name: "Error service invalid amount",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				svc := mock_transaction.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, entity.ErrInvalidAmount)

				return svc
			},
			reqBody:    []byte(`{"account_id": 1, "operation_type_id": 4, "amount": -123.45}`),
			wantStatus: http.StatusBadRequest,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{
					Title:  "Amount informed is Invalid",
					Detail: "invalid amount",
				})
			},
		},
		{
			name: "Error service insufficient credit limit",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				svc := mock_transaction.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, entity.ErrInsufficientCreditLimit)

				return svc
			},
			reqBody:    []byte(`{"account_id": 1, "operation_type_id": 4, "amount": -123.45}`),
			wantStatus: http.StatusBadRequest,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{
					Title:  "Insufficient Available Credit Limit",
					Detail: "available credit limit is insufficient",
				})
			},
		},
		{
			name: "Error service generic error",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				svc := mock_transaction.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error"))

				return svc
			},
			reqBody:    []byte(`{"account_id": 1, "operation_type_id": 1, "amount": 123.45}`),
			wantStatus: http.StatusInternalServerError,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{
					Title: "Error while creating Transaction",
				})
			},
		},
		{
			name: "Success",
			svcArgs: func(ctrl *gomock.Controller) transaction.Service {
				svc := mock_transaction.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, accountID, operationTypeID int, amount float64) (*entity.Transaction, error) {
						return &entity.Transaction{
							ID:              1,
							AccountID:       accountID,
							OperationTypeID: operationTypeID,
							Amount:          amount,
							EventDate:       time.Time{},
						}, nil
					})

				return svc
			},
			reqBody:    []byte(`{"account_id": 1, "operation_type_id": 4, "amount": 123.45}`),
			wantStatus: http.StatusCreated,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.TransactionResponse{
					ID:              1,
					AccountID:       1,
					OperationTypeID: 4,
					Amount:          123.45,
					EventDate:       time.Time{},
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()

			handler := NewTransactionHandler(tc.svcArgs(ctrl))

			app.Post("/transactions", handler.Create)

			wantBody, err := tc.wantBody()
			assert.NoError(t, err)

			apitest.New().
				HandlerFunc(testHelper.FiberToHandlerFunc(app)).
				Post("/transactions").
				Body(string(tc.reqBody)).
				Header(fiber.HeaderContentType, fiber.MIMEApplicationJSON).
				Expect(t).
				Status(tc.wantStatus).
				Body(string(wantBody)).
				End()
		})
	}
}
