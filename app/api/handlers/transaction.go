package handlers

import (
	"github.com/brunomdev/digital-account/app/api/presenter"
	"github.com/brunomdev/digital-account/domain/transaction"
	"github.com/brunomdev/digital-account/entity"
	"github.com/brunomdev/digital-account/infra/log"
	validator "github.com/brunomdev/digital-account/pkg/validate"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type TransactionHandler interface {
	Create(c *fiber.Ctx) error
}

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) TransactionHandler {
	return &transactionHandler{
		service: service,
	}
}

func (h *transactionHandler) Create(c *fiber.Ctx) error {
	var input struct {
		AccountID       int     `json:"account_id" validate:"required,min=1"`
		OperationTypeID int     `json:"operation_type_id" validate:"required,min=1"`
		Amount          float64 `json:"amount" validate:"required"`
	}

	err := c.BodyParser(&input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(
				presenter.ErrorResponse{
					Title:  "Unable to parse body",
					Detail: err.Error(),
				},
			)
	}

	errs := validator.ValidateStruct(input)
	if errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errs)
	}

	txn, err := h.service.Create(c.Context(), input.AccountID, input.OperationTypeID, input.Amount)
	if err != nil {
		log.Error(c.Context(), "unable to create transaction", err)

		errStatus, errResponse := h.formatErrResponse(err)

		return c.Status(errStatus).JSON(errResponse)
	}

	resp := presenter.TransactionResponse{
		ID:              txn.ID,
		AccountID:       txn.AccountID,
		OperationTypeID: txn.OperationTypeID,
		Amount:          txn.Amount,
		EventDate:       txn.EventDate,
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *transactionHandler) formatErrResponse(err error) (int, presenter.ErrorResponse) {
	errStatus := fiber.StatusInternalServerError
	errResponse := presenter.ErrorResponse{
		Title: "Error while creating Transaction",
	}

	switch {
	case errors.Is(err, entity.ErrNotFound):
		errStatus = fiber.StatusNotFound
		errResponse.Title = "Resource not found"
		errResponse.Detail = err.Error()
	case errors.Is(err, entity.ErrInvalidAmount):
		errStatus = fiber.StatusBadRequest
		errResponse.Title = "Amount informed is Invalid"
		errResponse.Detail = err.Error()
	case errors.Is(err, entity.ErrInsufficientCreditLimit):
		errStatus = fiber.StatusBadRequest
		errResponse.Title = "Insufficient Available Credit Limit"
		errResponse.Detail = err.Error()
	}

	return errStatus, errResponse
}
