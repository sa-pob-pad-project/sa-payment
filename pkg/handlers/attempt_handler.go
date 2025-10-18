package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// CreatePaymentAttempt godoc
// @Summary Create payment attempt
// @Description Create a new payment attempt for the authenticated patient
// @Tags payment-attempt
// @Accept json
// @Produce json
// @Param payment_attempt body dto.CreatePaymentAttemptRequestDto true "Payment attempt data"
// @Success 201 {object} dto.CreatePaymentAttemptResponseDto "Payment attempt created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or identifiers"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 404 {object} response.ErrorResponse "Payment information not found"
// @Failure 500 {object} response.ErrorResponse "Failed to create payment attempt"
// @Router /api/payment/v1/attempt [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) CreatePaymentAttempt(c *fiber.Ctx) error {
	var body dto.CreatePaymentAttemptRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	ctx := contextUtils.GetContext(c)

	res, err := h.paymentService.CreatePaymentAttempt(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return response.Created(c, res)
}

// GetPaymentAttempt godoc
// @Summary Get payment attempt by ID
// @Description Retrieve a payment attempt record by its identifier
// @Tags payment-attempt
// @Accept json
// @Produce json
// @Param id path string true "Payment attempt ID"
// @Success 200 {object} dto.GetPaymentAttemptResponseDto "Payment attempt retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid payment attempt ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Payment attempt not found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve payment attempt"
// @Router /api/payment/v1/attempt/{id} [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) GetPaymentAttempt(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing payment attempt ID",
		})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.GetPaymentAttempt(ctx, id)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
