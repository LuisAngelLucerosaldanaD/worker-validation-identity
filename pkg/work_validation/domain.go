package work_validation

import (
	"time"
)

// WorkValidation  Model struct WorkValidation
type WorkValidation struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	Status    string    `json:"status" db:"status" valid:"required"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
