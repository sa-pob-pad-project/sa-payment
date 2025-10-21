package service

import (
	"context"
	"encoding/json"
	"errors"
	"payment-service/pkg/apperr"
	"payment-service/pkg/clients"
	contextUtils "payment-service/pkg/context"
	"payment-service/pkg/dto"
	"payment-service/pkg/models"
	"payment-service/pkg/repository"
	"payment-service/pkg/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentService struct {
	db                           *gorm.DB
	paymentInformationRepository *repository.PaymentInformationRepository
	paymentAttemptRepository     *repository.PaymentAttemptRepository
	paymentRepository            *repository.PaymentRepository
	userClient                   *clients.UserClient
}

func NewPaymentService(
	db *gorm.DB,
	paymentInformationRepository *repository.PaymentInformationRepository,
	paymentAttemptRepository *repository.PaymentAttemptRepository,
	paymentRepository *repository.PaymentRepository,
	userClient *clients.UserClient,
) *PaymentService {
	return &PaymentService{
		db:                           db,
		paymentInformationRepository: paymentInformationRepository,
		paymentAttemptRepository:     paymentAttemptRepository,
		paymentRepository:            paymentRepository,
		userClient:                   userClient,
	}
}

// CreatePaymentInformation creates a new payment information record
func (s *PaymentService) CreatePaymentInfo(ctx context.Context, body dto.CreatePaymentInfoRequestDto) (*dto.CreatePaymentInfoResponseDto, error) {
	patientID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "patient" {
		return nil, apperr.New(apperr.CodeForbidden, "only patients can create payment information", nil)
	}

	detailsJSON, err := json.Marshal(body.Details)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment details", err)
	}

	paymentInfo := &models.PaymentInformation{
		ID:      utils.GenerateUUIDv7(),
		UserID:  utils.StringToUUIDv7(patientID),
		Type:    body.PaymentMethod,
		Details: detailsJSON,
		Version: 1,
	}

	if err := s.paymentInformationRepository.Create(ctx, paymentInfo); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to create payment info", err)
	}

	return &dto.CreatePaymentInfoResponseDto{
		PaymentInfo: dto.ToPaymentInfoDto(paymentInfo),
	}, nil
}

// GetPaymentInformationByID retrieves a payment information by ID
func (s *PaymentService) GetPaymentInfoByID(ctx context.Context, id string) (*dto.GetPaymentInfoByIDResponseDto, error) {
	paymentInfoID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment information ID format", err)
	}

	paymentInfo, err := s.paymentInformationRepository.FindByID(ctx, paymentInfoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment information not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment information", err)
	}

	return &dto.GetPaymentInfoByIDResponseDto{
		PaymentInfo: dto.ToPaymentInfoDto(paymentInfo),
	}, nil
}

func (s *PaymentService) GetPaymentInfoByMethod(ctx context.Context, method string) (*dto.GetPaymentInfoByIDResponseDto, error) {
	userID := contextUtils.GetUserId(ctx)
	paymentInfos, err := s.paymentInformationRepository.FindByUserIDAndType(ctx, utils.StringToUUIDv7(userID), method)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment information not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment information", err)
	}

	if len(paymentInfos) == 0 {
		return nil, apperr.New(apperr.CodeNotFound, "payment information not found", nil)
	}

	// For simplicity, return the first matching payment info
	return &dto.GetPaymentInfoByIDResponseDto{
		PaymentInfo: dto.ToPaymentInfoDto(&paymentInfos[0]),
	}, nil
}

// GetPaymentInformationByUserID retrieves all payment information for a user
func (s *PaymentService) GetAllPaymentInfos(ctx context.Context) (*dto.GetAllPaymentInfosResponseDto, error) {
	paymentInfos, err := s.paymentInformationRepository.FindAll(ctx)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve payment information", err)
	}

	return &dto.GetAllPaymentInfosResponseDto{
		DeliveryInfos: dto.ToPaymentInfoList(paymentInfos),
	}, nil
}

// UpdatePaymentInfo updates a payment information record
func (s *PaymentService) UpdatePaymentInfo(ctx context.Context, body dto.UpdatePaymentInfoRequestDto) (*dto.UpdatePaymentInfoResponseDto, error) {
	paymentID, err := uuid.Parse(body.ID)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment ID format", err)
	}

	// Get existing payment info
	existingPaymentInfo, err := s.paymentInformationRepository.FindByID(ctx, paymentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment information not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "Failed to retrieve payment information", err)
	}

	// Update fields
	if body.PaymentMethod != "" {
		existingPaymentInfo.Type = body.PaymentMethod
	}
	if body.Details != nil {
		existingPaymentInfo.Details = body.Details
	}
	existingPaymentInfo.Version += 1

	// Save updates
	err = s.paymentInformationRepository.Update(ctx, existingPaymentInfo)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "Failed to update payment information", err)
	}

	return &dto.UpdatePaymentInfoResponseDto{
		ID:            existingPaymentInfo.ID.String(),
		UserID:        existingPaymentInfo.UserID.String(),
		PaymentMethod: existingPaymentInfo.Type,
		Details:       existingPaymentInfo.Details,
		Version:       existingPaymentInfo.Version,
		CreatedAt:     existingPaymentInfo.CreatedAt,
	}, nil
}

// DeletePaymentInfo deletes a payment information record
func (s *PaymentService) DeletePaymentInfo(ctx context.Context, id string) (*dto.DeletePaymentInfoResponseDto, error) {
	paymentInfoID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperr.New(apperr.CodeBadRequest, "invalid payment information ID format", err)
	}

	_, err = s.paymentInformationRepository.FindByID(ctx, paymentInfoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.New(apperr.CodeNotFound, "payment information not found", nil)
		}
		return nil, apperr.New(apperr.CodeInternal, "failed to retrieve payment information", err)
	}

	err = s.paymentInformationRepository.Delete(ctx, paymentInfoID)
	if err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to delete payment information", err)
	}

	return &dto.DeletePaymentInfoResponseDto{
		ID:        id,
		DeletedAt: time.Now().Format(time.RFC3339),
	}, nil
}
