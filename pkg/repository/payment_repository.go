package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository) Transaction(ctx context.Context, fn func(repo *PaymentRepository) (interface{}, error)) (interface{}, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	repoWithTx := r.withTx(tx)

	result, err := fn(repoWithTx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PaymentRepository) withTx(tx *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: tx}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindByAttemptID(ctx context.Context, attemptID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.WithContext(ctx).Where("attempt_id = ?", attemptID).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) FindAll(ctx context.Context) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.db.WithContext(ctx).Find(&payments).Error; err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Model(payment).Updates(payment).Error
}

func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Payment{}).Error
}

func (r *PaymentRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&models.Payment{}).Error
}
