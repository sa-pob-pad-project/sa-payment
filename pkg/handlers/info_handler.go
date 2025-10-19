package handlers

import (
	"payment-service/pkg/apperr"
	contextUtils "payment-service/pkg/context"
	"payment-service/pkg/dto"
	"payment-service/pkg/response"
	service "payment-service/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePaymentInfo godoc
// @Summary Create payment information
// @Description Create a new payment method for the authenticated patient
// @Tags payment-info
// @Accept json
// @Produce json
// @Param payment_info body dto.CreatePaymentInfoRequestDto true "Payment information to create"
// @Success 201 {object} dto.CreatePaymentInfoResponseDto "Payment information created successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "Failed to create payment information"
// @Router /api/payment/v1/info [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) CreatePaymentInfo(c *fiber.Ctx) error {
	var body dto.CreatePaymentInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	ctx := contextUtils.GetContext(c)

	res, err := h.paymentService.CreatePaymentInfo(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return response.Created(c, res)
}

// GetPaymentInfo godoc
// @Summary Get payment information by ID
// @Description Retrieve a payment information record by its identifier
// @Tags payment-info
// @Accept json
// @Produce json
// @Param id path string true "Payment information ID"
// @Success 200 {object} dto.GetPaymentInfoByIDResponseDto "Payment information retrieved successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid payment information ID"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Payment information not found"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve payment information"
// @Router /api/payment/v1/info/{id} [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) GetPaymentInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing payment information ID",
		})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.GetPaymentInfoByID(ctx, id)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// GetAllPaymentInfos godoc
// @Summary List payment information
// @Description Retrieve all payment information records for the authenticated patient
// @Tags payment-info
// @Accept json
// @Produce json
// @Success 200 {object} dto.GetAllPaymentInfosResponseDto "Payment information retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to retrieve payment information"
// @Router /api/payment/v1/info [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) GetAllPaymentInfos(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.GetAllPaymentInfos(ctx)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// UpdatePaymentInfo godoc
// @Summary Update payment information
// @Description Update an existing payment information record
// @Tags payment-info
// @Accept json
// @Produce json
// @Param payment_info body dto.UpdatePaymentInfoRequestDto true "Payment information to update"
// @Success 200 {object} dto.UpdatePaymentInfoResponseDto "Payment information updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Payment information not found"
// @Failure 500 {object} response.ErrorResponse "Failed to update payment information"
// @Router /api/payment/v1/info [put]
// @Security ApiKeyAuth
func (h *PaymentHandler) UpdatePaymentInfo(c *fiber.Ctx) error {
	var body dto.UpdatePaymentInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body " + err.Error(),
		})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.UpdatePaymentInfo(ctx, body)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

// DeletePaymentInfo godoc
// @Summary Delete payment information
// @Description Delete an existing payment information record
// @Tags payment-info
// @Accept json
// @Produce json
// @Param payment_info body dto.DeletePaymentInfoRequestDto true "Payment information ID to delete"
// @Success 200 {object} dto.DeletePaymentInfoResponseDto "Payment information deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 404 {object} response.ErrorResponse "Payment information not found"
// @Failure 500 {object} response.ErrorResponse "Failed to delete payment information"
// @Router /api/payment/v1/info [delete]
// @Security ApiKeyAuth
func (h *PaymentHandler) DeletePaymentInfo(c *fiber.Ctx) error {
	var body dto.DeletePaymentInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body " + err.Error(),
		})
	}

	ctx := contextUtils.GetContext(c)
	res, err := h.paymentService.DeletePaymentInfo(ctx, body.ID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
