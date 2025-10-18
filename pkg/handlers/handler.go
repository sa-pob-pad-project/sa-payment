package handlers

import (
	service "order-service/pkg/services"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

// Example handler methods would go here

// func (h *PaymentHandler) CreatePaymentInformation(c fiber.Ctx) error {
// 	var body dto.Bodytype
// 	if err := c.BodyParser(&body); err != nil {
// 		return response.BadRequest(c, "Invalid request body "+err.Error())
// 	}
// 	ctx := contextUtils.GetContext(c)
// 	// Call service methods using h.paymentService
// 	res, err := h.paymentService.CreatePaymentInformation(ctx, &body)
// 	if err != nil {
// 		return apperr.WriteError(c, err)
// 	}
// 	return response.Created(c, res)
// }
