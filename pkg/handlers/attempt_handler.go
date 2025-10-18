package handlers

import (
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/response"

	"github.com/gofiber/fiber/v2"
)

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
