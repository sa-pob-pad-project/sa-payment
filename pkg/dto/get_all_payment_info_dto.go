package dto

import "payment-service/pkg/models"

type GetAllPaymentInfosResponseDto struct {
	DeliveryInfos []PaymentInfoDto `json:"delivery_infos"`
}

func ToPaymentInfoList(info []models.PaymentInformation) []PaymentInfoDto {
	result := make([]PaymentInfoDto, len(info))
	for i, v := range info {
		result[i] = ToPaymentInfoDto(&v)
	}
	return result
}
