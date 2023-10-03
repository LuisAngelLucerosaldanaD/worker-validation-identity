package onboarding

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Onboarding  Model struct Onboarding
type Onboarding struct {
	ID            string    `json:"id" db:"id" valid:"required,uuid"`
	ClientId      int64     `json:"client_id" db:"client_id" valid:"-"`
	RequestId     string    `json:"request_id" db:"request_id" valid:"required"`
	UserId        string    `json:"user_id" db:"user_id" valid:"required"`
	Status        string    `json:"status" db:"status" valid:"required"`
	TransactionId string    `json:"transaction_id" db:"transaction_id" valid:"-"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

func NewOnboarding(id string, clientId int64, requestId string, userId string, status string, transactionId string) *Onboarding {
	return &Onboarding{
		ID:            id,
		ClientId:      clientId,
		RequestId:     requestId,
		UserId:        userId,
		Status:        status,
		TransactionId: transactionId,
	}
}

func (m *Onboarding) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
