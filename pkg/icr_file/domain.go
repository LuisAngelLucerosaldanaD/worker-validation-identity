package icr_file

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// IcrFile  Model struct IcrFile
type IcrFile struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	Text      string    `json:"text" db:"text" valid:"required"`
	X         float64   `json:"x" db:"x" valid:"required"`
	Y         float64   `json:"y" db:"y" valid:"required"`
	W         float64   `json:"w" db:"w" valid:"required"`
	H         float64   `json:"h" db:"h" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewIcrFile(id int64, userId string, text string, x float64, y float64, w float64, h float64) *IcrFile {
	return &IcrFile{
		ID:     id,
		UserId: userId,
		Text:   text,
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
	}
}

func NewCreateIcrFile(userId string, text string, x float64, y float64, w float64, h float64) *IcrFile {
	return &IcrFile{
		UserId: userId,
		Text:   text,
		X:      x,
		Y:      y,
		W:      w,
		H:      h,
	}
}

func (m *IcrFile) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
