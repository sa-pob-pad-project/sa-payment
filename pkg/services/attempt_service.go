package service

import (
	"context"
	"errors"
	"order-service/pkg/apperr"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/models"
	"order-service/pkg/utils"

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
