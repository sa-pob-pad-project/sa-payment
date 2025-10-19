package repository

import (
	"context"
	"payment-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentAttemptRepository struct {
	db *gorm.DB
}

func NewPaymentAttemptRepository(db *gorm.DB) *PaymentAttemptRepository {
	return &PaymentAttemptRepository{
		db: db,
	}
}

func (r *PaymentAttemptRepository) Transaction(ctx context.Context, fn func(repo *PaymentAttemptRepository) (interface{}, error)) (interface{}, error) {
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

func (r *PaymentAttemptRepository) withTx(tx *gorm.DB) *PaymentAttemptRepository {
	return &PaymentAttemptRepository{db: tx}
}

func (r *PaymentAttemptRepository) Create(ctx context.Context, attempt *models.PaymentAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}

func (r *PaymentAttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PaymentAttempt, error) {
	var attempt models.PaymentAttempt
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&attempt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &attempt, nil
}

func (r *PaymentAttemptRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.PaymentAttempt, error) {
	var attempts []models.PaymentAttempt
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *PaymentAttemptRepository) FindByOrderIDAndStatus(ctx context.Context, orderID uuid.UUID, status models.PaymentStatus) ([]models.PaymentAttempt, error) {
	var attempts []models.PaymentAttempt
	if err := r.db.WithContext(ctx).Where("order_id = ? AND status = ?", orderID, status).Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *PaymentAttemptRepository) FindAll(ctx context.Context) ([]models.PaymentAttempt, error) {
	var attempts []models.PaymentAttempt
	if err := r.db.WithContext(ctx).Find(&attempts).Error; err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *PaymentAttemptRepository) Update(ctx context.Context, attempt *models.PaymentAttempt) error {
	return r.db.WithContext(ctx).Model(attempt).Updates(attempt).Error
}

func (r *PaymentAttemptRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PaymentAttempt{}).Error
}

func (r *PaymentAttemptRepository) DeleteByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&models.PaymentAttempt{}).Error
}
