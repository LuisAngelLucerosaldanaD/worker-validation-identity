package file

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// File  Model struct File
type File struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	Path      string    `json:"path" db:"path" valid:"required"`
	Name      string    `json:"name" db:"name" valid:"required"`
	Type      int32     `json:"type" db:"type" valid:"required"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewFile(id int64, path string, name string, typeFile int32, userId string) *File {
	return &File{
		ID:     id,
		Path:   path,
		Name:   name,
		Type:   typeFile,
		UserId: userId,
	}
}

func NewCreateFile(path string, name string, typeFile int32, userId string) *File {
	return &File{
		Path:   path,
		Name:   name,
		Type:   typeFile,
		UserId: userId,
	}
}

func (m *File) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
