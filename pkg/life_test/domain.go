package life_test

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// LifeTest  Model struct LifeTest
type LifeTest struct {
	ID         int64     `json:"id" db:"id" valid:"-"`
	ClientId   int64     `json:"client_id" db:"client_id" valid:"required"`
	MaxNumTest int       `json:"max_num_test" db:"max_num_test" valid:"required"`
	RequestId  string    `json:"request_id" db:"request_id" valid:"required"`
	ExpiredAt  time.Time `json:"expired_at" db:"expired_at" valid:"required"`
	UserID     string    `json:"user_id" db:"user_id" valid:"required"`
	Status     string    `json:"status" db:"status" valid:"required"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

func NewLifeTest(id int64, clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) *LifeTest {
	return &LifeTest{
		ID:         id,
		ClientId:   clientId,
		MaxNumTest: maxNumValidation,
		RequestId:  requestId,
		ExpiredAt:  expiredAt,
		UserID:     userIdentification,
		Status:     status,
	}
}

func NewCreateLifeTest(clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) *LifeTest {
	return &LifeTest{
		ClientId:   clientId,
		MaxNumTest: maxNumValidation,
		RequestId:  requestId,
		ExpiredAt:  expiredAt,
		UserID:     userIdentification,
		Status:     status,
	}
}

func (m *LifeTest) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
