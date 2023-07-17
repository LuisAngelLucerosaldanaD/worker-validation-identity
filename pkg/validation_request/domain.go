package validation_request

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// ValidationRequest  Model struct ValidationRequest
type ValidationRequest struct {
	ID                 int64     `json:"id" db:"id" valid:"-"`
	ClientId           int64     `json:"client_id" db:"client_id" valid:"required"`
	MaxNumValidation   int       `json:"max_num_validation" db:"max_num_validation" valid:"required"`
	RequestId          string    `json:"request_id" db:"request_id" valid:"required"`
	ExpiredAt          time.Time `json:"expired_at" db:"expired_at" valid:"required"`
	UserIdentification string    `json:"user_identification" db:"user_identification" valid:"required"`
	Status             string    `json:"status" db:"status" valid:"required"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

func NewValidationRequest(id int64, clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) *ValidationRequest {
	return &ValidationRequest{
		ID:                 id,
		ClientId:           clientId,
		MaxNumValidation:   maxNumValidation,
		RequestId:          requestId,
		ExpiredAt:          expiredAt,
		UserIdentification: userIdentification,
		Status:             status,
	}
}

func NewCreateValidationRequest(clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) *ValidationRequest {
	return &ValidationRequest{
		ClientId:           clientId,
		MaxNumValidation:   maxNumValidation,
		RequestId:          requestId,
		ExpiredAt:          expiredAt,
		UserIdentification: userIdentification,
		Status:             status,
	}
}

func (m *ValidationRequest) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
