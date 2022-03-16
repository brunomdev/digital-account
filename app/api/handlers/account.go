package handlers

import (
	"github.com/brunomdev/digital-account/app/api/presenter"
	"github.com/brunomdev/digital-account/domain/account"
	"github.com/brunomdev/digital-account/infra/log"
	validator "github.com/brunomdev/digital-account/pkg/validate"
	"github.com/gofiber/fiber/v2"
)

type AccountHandler interface {
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
}

type accountHandler struct {
	service account.Service
}

func NewAccountHandler(service account.Service) AccountHandler {
	return &accountHandler{
		service: service,
	}
}

func (h *accountHandler) Create(c *fiber.Ctx) error {
	var input struct {
		DocumentNumber string `json:"document_number" validate:"required"`
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

	acc, err := h.service.Create(c.Context(), input.DocumentNumber)
	if err != nil {
		log.Error(c.Context(), "unable to create account", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			presenter.ErrorResponse{Title: "Error while creating Account"},
		)
	}

	resp := presenter.AccountResponse{
		ID:             acc.ID,
		DocumentNumber: acc.DocumentNumber,
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *accountHandler) Get(c *fiber.Ctx) error {
	var input struct {
		ID int `validate:"required,min=1"`
	}

	input.ID, _ = c.ParamsInt("id")

	errs := validator.ValidateStruct(input)
	if errs != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errs)
	}

	acc, err := h.service.Get(c.Context(), input.ID)
	if err != nil {
		log.Error(c.Context(), "unable to find account", err)

		return c.Status(fiber.StatusInternalServerError).JSON(
			presenter.ErrorResponse{Title: "Account not found"},
		)
	}

	resp := presenter.AccountResponse{
		ID:             acc.ID,
		DocumentNumber: acc.DocumentNumber,
	}

	return c.JSON(resp)
}
