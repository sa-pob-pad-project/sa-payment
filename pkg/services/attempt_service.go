package service

import (
	"context"
	"errors"
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/models"
	"order-service/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *PaymentService) CreatePaymentAttempt(ctx context.Context, body dto.CreatePaymentAttemptRequestDto) (*dto.CreatePaymentAttemptResponseDto, error) {
	role := contextUtils.GetRole(ctx)
	if role != "patient" {
		return nil, apperr.New(apperr.CodeForbidden, "only patients can create payment attempts", nil)
	}

	orderID := utils.StringToUUIDv7(body.OrderID)
	if orderID == uuid.Nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid order ID", nil)
	}

	paymentInfoID := utils.StringToUUIDv7(body.PaymentInfoID)
	if paymentInfoID == uuid.Nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment information ID", nil)
	}

	paymentInfo, err := s.paymentInformationRepository.FindByID(ctx, paymentInfoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment information not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment information", err)
	}

	paymentAttempt := &models.PaymentAttempt{
		ID:                   utils.GenerateUUIDv7(),
		OrderID:              orderID,
		PaymentInformationID: &paymentInfoID,
		Method:               paymentInfo.Type,
		Status:               models.PaymentStatusPending,
	}

	if err := s.paymentAttemptRepository.Create(ctx, paymentAttempt); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to create payment attempt", err)
	}

	return &dto.CreatePaymentAttemptResponseDto{
		PaymentAttemptID: paymentAttempt.ID.String(),
	}, nil
}

func (s *PaymentService) GetPaymentAttempt(ctx context.Context, paymentAttemptID string) (*dto.GetPaymentAttemptResponseDto, error) {
	id := utils.StringToUUIDv7(paymentAttemptID)
	if id == uuid.Nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment attempt ID", nil)
	}

	paymentAttempt, err := s.paymentAttemptRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment attempt not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment attempt", err)
	}

	response := &dto.GetPaymentAttemptResponseDto{
		PaymentAttemptID: paymentAttempt.ID.String(),
		OrderID:          paymentAttempt.OrderID.String(),
		Method:           paymentAttempt.Method,
		Status:           paymentAttempt.Status,
	}

	if paymentAttempt.PaymentInformationID != nil {
		response.PaymentInfoID = paymentAttempt.PaymentInformationID.String()
	}

	return response, nil
}

func (s *PaymentService) UpdatePaymentAttempt(ctx context.Context, body dto.UpdatePaymentAttemptRequestDto) (*dto.UpdatePaymentAttemptResponseDto, error) {
	if body.Status == "" {
		return nil, apperr.New(apperr.CodeBadRequest, "status is required", nil)
	}

	if !isValidPaymentStatus(body.Status) {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment status provided", nil)
	}

	id := utils.StringToUUIDv7(body.PaymentAttemptID)
	if id == uuid.Nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment attempt ID", nil)
	}

	paymentAttempt, err := s.paymentAttemptRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment attempt not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment attempt", err)
	}

	paymentAttempt.Status = body.Status

	if err := s.paymentAttemptRepository.Update(ctx, paymentAttempt); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to update payment attempt", err)
	}

	response := &dto.UpdatePaymentAttemptResponseDto{
		PaymentAttemptID: paymentAttempt.ID.String(),
		OrderID:          paymentAttempt.OrderID.String(),
		Method:           paymentAttempt.Method,
		Status:           paymentAttempt.Status,
	}

	if paymentAttempt.PaymentInformationID != nil {
		response.PaymentInfoID = paymentAttempt.PaymentInformationID.String()
	}

	return response, nil
}

func isValidPaymentStatus(status models.PaymentStatus) bool {
	switch status {
	case models.PaymentStatusPending,
		models.PaymentStatusSuccess,
		models.PaymentStatusFailed:
		return true
	default:
		return false
	}
}

func (s *PaymentService) CreatePayment(ctx context.Context, body dto.CreatePaymentRequestDto) (*dto.CreatePaymentResponseDto, error) {
	if body.PaymentAttemptID == "" {
		return nil, apperr.New(apperr.CodeBadRequest, "attempt ID is required", nil)
	}

	if body.Amount < 0 {
		return nil, apperr.New(apperr.CodeBadRequest, "amount must be greater than or equal to zero", nil)
	}

	attemptID := utils.StringToUUIDv7(body.PaymentAttemptID)
	if attemptID == uuid.Nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid attempt ID", nil)
	}

	paymentAttempt, err := s.paymentAttemptRepository.FindByID(ctx, attemptID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment attempt not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment attempt", err)
	}

	if paymentAttempt.Status != models.PaymentStatusSuccess {
		return nil, apperr.New(apperr.CodeBadRequest, "payment can only be created for successful attempts", nil)
	}

	payment := &models.Payment{
		ID:        utils.GenerateUUIDv7(),
		AttemptID: attemptID,
		Amount:    body.Amount,
		OrderID:   paymentAttempt.OrderID,
		PaidAt:    time.Now().UTC(),
	}

	if err := s.paymentRepository.Create(ctx, payment); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to create payment", err)
	}

	return &dto.CreatePaymentResponseDto{
		PaymentID: payment.ID.String(),
		AttemptID: payment.AttemptID.String(),
		OrderID:   payment.OrderID.String(),
		Amount:    payment.Amount,
		PaidAt:    payment.PaidAt.Format(time.RFC3339),
	}, nil
}
