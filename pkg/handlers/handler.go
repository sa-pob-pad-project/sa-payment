package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/response"
	service "order-service/pkg/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// CreatePaymentInfo creates a new payment information record
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

// GetPaymentInfo retrieves a payment information by ID
func (h *PaymentHandler) GetPaymentInfo(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := contextUtils.GetContext(c)

	paymentID, err := uuid.Parse(id)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID format")
	}

	paymentInfo, err := h.paymentService.GetPaymentInformationByID(ctx, paymentID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	res := &dto.GetPaymentInfoResponseDto{
		ID:            paymentInfo.ID.String(),
		UserID:        paymentInfo.UserID.String(),
		PaymentMethod: paymentInfo.Type,
		Details:       paymentInfo.Details,
		Version:       paymentInfo.Version,
		CreatedAt:     paymentInfo.CreatedAt,
	}

	return response.OK(c, res)
}

// GetAllPaymentInfo retrieves all payment information for the authenticated user
func (h *PaymentHandler) GetAllPaymentInfo(c *fiber.Ctx) error {
	ctx := contextUtils.GetContext(c)

	userIDStr := contextUtils.GetUserId(ctx)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "Invalid user ID format")
	}

	paymentInfos, err := h.paymentService.GetAllPaymentInfo(ctx, userID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	var res []dto.GetPaymentInfoResponseDto
	for _, pi := range paymentInfos {
		res = append(res, dto.GetPaymentInfoResponseDto{
			ID:            pi.ID.String(),
			UserID:        pi.UserID.String(),
			PaymentMethod: pi.Type,
			Details:       pi.Details,
			Version:       pi.Version,
			CreatedAt:     pi.CreatedAt,
		})
	}

	return response.OK(c, res)
}

// UpdatePaymentInfo updates a payment information record
func (h *PaymentHandler) UpdatePaymentInfo(c *fiber.Ctx) error {
	var body dto.UpdatePaymentInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	ctx := contextUtils.GetContext(c)

	paymentID, err := uuid.Parse(body.ID)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID format")
	}

	// Get existing payment info
	existingPaymentInfo, err := h.paymentService.GetPaymentInformationByID(ctx, paymentID)
	if err != nil {
		return apperr.WriteError(c, err)
	}

	// Update fields
	if body.PaymentMethod != "" {
		existingPaymentInfo.Type = body.PaymentMethod
	}
	if body.Details != nil {
		existingPaymentInfo.Details = body.Details
	}
	if body.Version > 0 {
		existingPaymentInfo.Version = body.Version
	}

	if err := h.paymentService.UpdatePaymentInformation(ctx, existingPaymentInfo); err != nil {
		return apperr.WriteError(c, err)
	}

	res := &dto.GetPaymentInfoResponseDto{
		ID:            existingPaymentInfo.ID.String(),
		UserID:        existingPaymentInfo.UserID.String(),
		PaymentMethod: existingPaymentInfo.Type,
		Details:       existingPaymentInfo.Details,
		Version:       existingPaymentInfo.Version,
		CreatedAt:     existingPaymentInfo.CreatedAt,
	}

	return response.OK(c, res)
}

// DeletePaymentInfo deletes a payment information record
func (h *PaymentHandler) DeletePaymentInfo(c *fiber.Ctx) error {
	var body dto.DeletePaymentInfoRequestDto
	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "Invalid request body "+err.Error())
	}
	ctx := contextUtils.GetContext(c)

	paymentID, err := uuid.Parse(body.ID)
	if err != nil {
		return response.BadRequest(c, "Invalid payment ID format")
	}

	if err := h.paymentService.DeletePaymentInformation(ctx, paymentID); err != nil {
		return apperr.WriteError(c, err)
	}

	return response.OK(c, map[string]string{"message": "Payment information deleted successfully"})
}
