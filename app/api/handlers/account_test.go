package handlers

import (
	"context"
	"encoding/json"
	"github.com/brunomdev/digital-account/app/api/presenter"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/domain/account/mock_account"
	"github.com/brunomdev/digital-account/entity"
	testHelper "github.com/brunomdev/digital-account/pkg/tests"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_accountHandler_Create(t *testing.T) {
	testCases := []struct {
		name       string
		svcArgs    func(ctrl *gomock.Controller) account.Service
		reqBody    []byte
		wantStatus int
		wantBody   func() ([]byte, error)
	}{
		{
			name: "Error bodyParser",
			svcArgs: func(ctrl *gomock.Controller) account.Service {
				return mock_account.NewMockService(ctrl)
			},
			reqBody:    []byte(`{"document_number": "12345678900",}`),
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
			svcArgs: func(ctrl *gomock.Controller) account.Service {
				return mock_account.NewMockService(ctrl)
			},
			reqBody:    []byte(`{}`),
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: func() ([]byte, error) {
				return json.Marshal([]presenter.ErrorResponse{
					{
						Source: "DocumentNumber",
						Detail: "DocumentNumber is a required field",
					},
				})
			},
		},
		{
			name: "Error service",
			svcArgs: func(ctrl *gomock.Controller) account.Service {
				svc := mock_account.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

				return svc
			},
			reqBody:    []byte(`{"document_number": "12345678900"}`),
			wantStatus: http.StatusInternalServerError,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{Title: "Error while creating Account"})
			},
		},
		{
			name: "Success",
			svcArgs: func(ctrl *gomock.Controller) account.Service {
				svc := mock_account.NewMockService(ctrl)

				svc.EXPECT().Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, docNumber string) (*entity.Account, error) {
						return &entity.Account{
							ID:             1,
							DocumentNumber: docNumber,
						}, nil
					})

				return svc
			},
			reqBody:    []byte(`{"document_number": "12345678900"}`),
			wantStatus: http.StatusCreated,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.AccountResponse{
					ID:             1,
					DocumentNumber: "12345678900",
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()

			handler := NewAccountHandler(tc.svcArgs(ctrl))

			app.Post("/accounts", handler.Create)

			wantBody, err := tc.wantBody()
			assert.NoError(t, err)

			apitest.New().
				HandlerFunc(testHelper.FiberToHandlerFunc(app)).
				Post("/accounts").
				Body(string(tc.reqBody)).
				Header(fiber.HeaderContentType, fiber.MIMEApplicationJSON).
				Expect(t).
				Status(tc.wantStatus).
				Body(string(wantBody)).
				End()
		})
	}
}

func Test_accountHandler_Get(t *testing.T) {
	type args struct {
		accountID int
	}
	testCases := []struct {
		name         string
		eventService func(ctrl *gomock.Controller) account.Service
		args         args
		wantStatus   int
		wantBody     func() ([]byte, error)
	}{
		{
			name: "Error invalid ID",
			eventService: func(ctrl *gomock.Controller) account.Service {
				return mock_account.NewMockService(ctrl)
			},
			args:       args{accountID: 0},
			wantStatus: http.StatusUnprocessableEntity,
			wantBody: func() ([]byte, error) {
				return json.Marshal([]presenter.ErrorResponse{
					{
						Source: "ID",
						Detail: "ID is a required field",
					},
				})
			},
		},
		{
			name: "Error service",
			eventService: func(ctrl *gomock.Controller) account.Service {
				svc := mock_account.NewMockService(ctrl)

				svc.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

				return svc
			},
			args:       args{accountID: 1},
			wantStatus: http.StatusInternalServerError,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.ErrorResponse{
					Title: "Account not found",
				})
			},
		},
		{
			name: "Error service",
			eventService: func(ctrl *gomock.Controller) account.Service {
				svc := mock_account.NewMockService(ctrl)

				svc.EXPECT().Get(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, id int) (*entity.Account, error) {
						return &entity.Account{
							ID:             id,
							DocumentNumber: "12345678900",
						}, nil
					})

				return svc
			},
			args:       args{accountID: 2},
			wantStatus: http.StatusOK,
			wantBody: func() ([]byte, error) {
				return json.Marshal(presenter.AccountResponse{
					ID:             2,
					DocumentNumber: "12345678900",
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			app := fiber.New()

			handler := NewAccountHandler(tc.eventService(ctrl))

			app.Get("/accounts/:id", handler.Get)

			wantBody, err := tc.wantBody()
			assert.NoError(t, err)

			apitest.New().
				HandlerFunc(testHelper.FiberToHandlerFunc(app)).
				Getf("/accounts/%d", tc.args.accountID).
				Header(fiber.HeaderContentType, fiber.MIMEApplicationJSON).
				Expect(t).
				Status(tc.wantStatus).
				Body(string(wantBody)).
				End()
		})
	}
}
