package traceability

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Traceability  Model struct Traceability
type Traceability struct {
	ID          int64     `json:"id" db:"id" valid:"-"`
	Action      string    `json:"action" db:"action" valid:"required"`
	Type        string    `json:"type" db:"type" valid:"required"`
	Description string    `json:"description" db:"description" valid:"required"`
	UserId      string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewTraceability(id int64, action string, typeTrx string, description string, userId string) *Traceability {
	return &Traceability{
		ID:          id,
		Action:      action,
		Type:        typeTrx,
		Description: description,
		UserId:      userId,
	}
}

func NewCreateTraceability(action string, typeTrx string, description string, userId string) *Traceability {
	return &Traceability{
		Action:      action,
		Type:        typeTrx,
		Description: description,
		UserId:      userId,
	}
}

func (m *Traceability) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
