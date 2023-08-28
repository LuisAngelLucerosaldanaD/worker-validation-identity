package validation_identity

import "time"

type RequestOnboarding struct {
	DocumentNumber string    `json:"document_number"`
	Status         string    `json:"status"`
	RequestID      string    `json:"request_id"`
	UserID         string    `json:"user_id"`
	VerifiedAt     time.Time `json:"verified_at"`
}
