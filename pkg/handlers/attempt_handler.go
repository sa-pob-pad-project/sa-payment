package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// CreatePayment godoc
// @Summary Create payment
// @Description Create a payment record for a successful payment attempt
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body dto.CreatePaymentRequestDto true "Payment creation payload"
// @Success 201 {object} dto.CreatePaymentResponseDto "Payment created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or business rule violation"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Payment attempt not found"
// @Failure 500 {object} response.ErrorResponse "Failed to create payment"
// @Router /api/payment/v1/ [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) CreatePayment(c *fiber.Ctx) error {
	var body dto.CreatePaymentRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.CreatePayment(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return response.Created(c, res)
}

// GetAllPayments godoc
// @Summary List payments
// @Description Retrieve all payment records
// @Tags payments
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAllPaymentsResponseDto "Payments retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve payments"
// @Router /api/payment/v1/ [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) GetAllPayments(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.GetAllPayments(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return response.OK(c, res)
}

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

// UpdatePaymentAttempt godoc
// @Summary Update payment attempt status
// @Description Update the status of an existing payment attempt
// @Tags payment-attempt
// @Accept json
// @Produce json
// @Param payment_attempt body dto.UpdatePaymentAttemptRequestDto true "Payment attempt status update payload"
// @Success 200 {object} dto.UpdatePaymentAttemptResponseDto "Payment attempt updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body or identifiers"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Payment attempt not found"
// @Failure 500 {object} response.ErrorResponse "Failed to update payment attempt"
// @Router /api/payment/v1/attempt [patch]
// @Security ApiKeyAuth
func (h *PaymentHandler) UpdatePaymentAttempt(c *fiber.Ctx) error {
	var body dto.UpdatePaymentAttemptRequestDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body " + err.Error(),
		})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.UpdatePaymentAttempt(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
