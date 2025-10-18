package repository

import (
	"context"
	"order-service/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentInformationRepository struct {
	db *gorm.DB
}

func NewPaymentInformationRepository(db *gorm.DB) *PaymentInformationRepository {
	return &PaymentInformationRepository{
		db: db,
	}
}

func (r *PaymentInformationRepository) Transaction(ctx context.Context, fn func(repo *PaymentInformationRepository) (interface{}, error)) (interface{}, error) {
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

func (r *PaymentInformationRepository) withTx(tx *gorm.DB) *PaymentInformationRepository {
	return &PaymentInformationRepository{db: tx}
}

func (r *PaymentInformationRepository) Create(ctx context.Context, paymentInfo *models.PaymentInformation) error {
	return r.db.WithContext(ctx).Create(paymentInfo).Error
}

func (r *PaymentInformationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.PaymentInformation, error) {
	var paymentInfo models.PaymentInformation
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&paymentInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}
	return &paymentInfo, nil
}

func (r *PaymentInformationRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.PaymentInformation, error) {
	var paymentInfos []models.PaymentInformation
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&paymentInfos).Error; err != nil {
		return nil, err
	}
	return paymentInfos, nil
}

func (r *PaymentInformationRepository) FindByUserIDAndType(ctx context.Context, userID uuid.UUID, paymentMethod models.PaymentMethod) ([]models.PaymentInformation, error) {
	var paymentInfos []models.PaymentInformation
	if err := r.db.WithContext(ctx).Where("user_id = ? AND type = ?", userID, paymentMethod).Find(&paymentInfos).Error; err != nil {
		return nil, err
	}
	return paymentInfos, nil
}

func (r *PaymentInformationRepository) FindAll(ctx context.Context) ([]models.PaymentInformation, error) {
	var paymentInfos []models.PaymentInformation
	if err := r.db.WithContext(ctx).Find(&paymentInfos).Error; err != nil {
		return nil, err
	}
	return paymentInfos, nil
}

func (r *PaymentInformationRepository) Update(ctx context.Context, paymentInfo *models.PaymentInformation) error {
	return r.db.WithContext(ctx).Model(paymentInfo).Updates(paymentInfo).Error
}

func (r *PaymentInformationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PaymentInformation{}).Error
}

func (r *PaymentInformationRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.PaymentInformation{}).Error
}
