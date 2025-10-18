package service

import (
	"context"
	"order-service/pkg/apperr"
	"order-service/pkg/clients"
	contextUtils "order-service/pkg/context"
	"order-service/pkg/dto"
	"order-service/pkg/models"
	"order-service/pkg/repository"
	"order-service/pkg/utils"

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

// func (s *OrderService) CreateOrder(ctx context.Context, body dto.CreateOrderRequestDto) (*dto.CreateOrderResponseDto, error) {
// 	userID := contextUtils.GetUserId(ctx)
// 	role := contextUtils.GetRole(ctx)

// 	if role != "patient" {
// 		return nil, apperr.New(apperr.CodeForbidden, "only patients can create orders", nil)
// 	}
// 	patientID, err := uuid.Parse(userID)
// 	if err != nil {
// 		return nil, apperr.New(apperr.CodeBadRequest, "invalid user ID", err)
// 	}
// 	appointment, err := s.appointmentClient.GetLatestAppointmentByPatientID(ctx, patientID)
// 	if err != nil {
// 		return nil, apperr.New(apperr.CodeInternal, "failed to get latest appointment", err)
// 	}
// 	if appointment == nil {
// 		return nil, apperr.New(apperr.CodeBadRequest, "patient has no appointment history", nil)
// 	}

// 	did := utils.StringToUUIDv7(appointment.DoctorID)
// 	order := &models.Order{
// 		ID:        utils.GenerateUUIDv7(),
// 		PatientID: patientID,
// 		DoctorID:  &did,
// 		Note:      body.Note,
// 		Status:    models.OrderStatusPending,
// 	}

// 	if err := s.orderRepository.Create(ctx, order); err != nil {
// 		return nil, apperr.New(apperr.CodeInternal, "failed to create order", err)
// 	}

// 	return &dto.CreateOrderResponseDto{
// 		OrderID: order.ID.String(),
// 	}, nil
// }

// CreatePaymentInformation creates a new payment information record
func (s *PaymentService) CreatePaymentInfo(ctx context.Context, body dto.CreatePaymentInfoRequestDto) (*dto.CreatePaymentInfoResponseDto, error) {
	patientID := contextUtils.GetUserId(ctx)
	role := contextUtils.GetRole(ctx)

	if role != "patient" {
		return nil, apperr.New(apperr.CodeForbidden, "only patients can create payment information", nil)
	}

	paymentInfo := &models.PaymentInformation{
		ID:      utils.GenerateUUIDv7(),
		UserID:  utils.StringToUUIDv7(patientID),
		Type:    body.PaymentMethod,
		Details: body.Details,
		Version: 1,
	}

	if err := s.paymentInformationRepository.Create(ctx, paymentInfo); err != nil {
		return nil, apperr.New(apperr.CodeInternal, "failed to create payment info", err)
	}

	return &dto.CreatePaymentInfoResponseDto{
		ID:            paymentInfo.ID.String(),
		UserID:        paymentInfo.UserID.String(),
		PaymentMethod: paymentInfo.Type,
		Details:       paymentInfo.Details,
		Version:       paymentInfo.Version,
		CreatedAt:     paymentInfo.CreatedAt,
	}, nil
}

// GetPaymentInformationByID retrieves a payment information by ID
func (s *PaymentService) GetPaymentInformationByID(ctx context.Context, id uuid.UUID) (*models.PaymentInformation, error) {
	return s.paymentInformationRepository.FindByID(ctx, id)
}

// GetPaymentInformationByUserID retrieves all payment information for a user
func (s *PaymentService) GetAllPaymentInfo(ctx context.Context, userID uuid.UUID) ([]models.PaymentInformation, error) {
	return s.paymentInformationRepository.FindByUserID(ctx, userID)
}

// GetPaymentInformationByUserIDAndType retrieves payment information by user and method
func (s *PaymentService) GetPaymentInformationByUserIDAndType(ctx context.Context, userID uuid.UUID, method models.PaymentMethod) ([]models.PaymentInformation, error) {
	return s.paymentInformationRepository.FindByUserIDAndType(ctx, userID, method)
}

// UpdatePaymentInformation updates a payment information record
func (s *PaymentService) UpdatePaymentInformation(ctx context.Context, paymentInfo *models.PaymentInformation) error {
	return s.paymentInformationRepository.Update(ctx, paymentInfo)
}

// DeletePaymentInformation deletes a payment information record
func (s *PaymentService) DeletePaymentInformation(ctx context.Context, id uuid.UUID) error {
	return s.paymentInformationRepository.Delete(ctx, id)
}

// CreatePaymentAttempt creates a new payment attempt
func (s *PaymentService) CreatePaymentAttempt(ctx context.Context, attempt *models.PaymentAttempt) error {
	return s.paymentAttemptRepository.Create(ctx, attempt)
}

// GetPaymentAttemptByID retrieves a payment attempt by ID
func (s *PaymentService) GetPaymentAttemptByID(ctx context.Context, id uuid.UUID) (*models.PaymentAttempt, error) {
	return s.paymentAttemptRepository.FindByID(ctx, id)
}

// GetPaymentAttemptsByOrderID retrieves all payment attempts for an order
func (s *PaymentService) GetPaymentAttemptsByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.PaymentAttempt, error) {
	return s.paymentAttemptRepository.FindByOrderID(ctx, orderID)
}

// GetPaymentAttemptsByOrderIDAndStatus retrieves payment attempts by order and status
func (s *PaymentService) GetPaymentAttemptsByOrderIDAndStatus(ctx context.Context, orderID uuid.UUID, status models.PaymentStatus) ([]models.PaymentAttempt, error) {
	return s.paymentAttemptRepository.FindByOrderIDAndStatus(ctx, orderID, status)
}

// UpdatePaymentAttempt updates a payment attempt
func (s *PaymentService) UpdatePaymentAttempt(ctx context.Context, attempt *models.PaymentAttempt) error {
	return s.paymentAttemptRepository.Update(ctx, attempt)
}

// DeletePaymentAttempt deletes a payment attempt
func (s *PaymentService) DeletePaymentAttempt(ctx context.Context, id uuid.UUID) error {
	return s.paymentAttemptRepository.Delete(ctx, id)
}

// CreatePayment creates a new payment record
func (s *PaymentService) CreatePayment(ctx context.Context, payment *models.Payment) error {
	return s.paymentRepository.Create(ctx, payment)
}

// GetPaymentByID retrieves a payment by ID
func (s *PaymentService) GetPaymentByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	return s.paymentRepository.FindByID(ctx, id)
}

// GetPaymentsByOrderID retrieves all payments for an order
func (s *PaymentService) GetPaymentsByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepository.FindByOrderID(ctx, orderID)
}

// GetPaymentsByAttemptID retrieves all payments for an attempt
func (s *PaymentService) GetPaymentsByAttemptID(ctx context.Context, attemptID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepository.FindByAttemptID(ctx, attemptID)
}

// UpdatePayment updates a payment
func (s *PaymentService) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	return s.paymentRepository.Update(ctx, payment)
}

// DeletePayment deletes a payment
func (s *PaymentService) DeletePayment(ctx context.Context, id uuid.UUID) error {
	return s.paymentRepository.Delete(ctx, id)
}
