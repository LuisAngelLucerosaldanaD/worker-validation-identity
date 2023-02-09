package status_request

import (
	"time"
)

// StatusRequest  Model struct StatusRequest
type StatusRequest struct {
	ID          int64     `json:"id" db:"id" valid:"-"`
	Status      string    `json:"status" db:"status" valid:"required"`
	Description string    `json:"description" db:"description" valid:"required"`
	UserId      string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewStatusRequest(id int64, status string, description string, userId string) *StatusRequest {
	return &StatusRequest{
		ID:          id,
		Status:      status,
		Description: description,
		UserId:      userId,
	}
}

func NewCreateStatusRequest(status string, description string, userId string) *StatusRequest {
	return &StatusRequest{
		Status:      status,
		Description: description,
		UserId:      userId,
	}
}
