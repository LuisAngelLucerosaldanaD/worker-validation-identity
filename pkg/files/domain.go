package files

import (
	"time"
)

// Files  Model struct Files
type Files struct {
	ID        int64     `json:"id" db:"id" valid:"-"`
	Path      string    `json:"path" db:"path" valid:"required"`
	Name      string    `json:"name" db:"name" valid:"required"`
	Type      int32     `json:"type" db:"type" valid:"required"`
	UserId    string    `json:"user_id" db:"user_id" valid:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewFiles(id int64, path string, name string, typeFile int32, userId string) *Files {
	return &Files{
		ID:     id,
		Path:   path,
		Name:   name,
		Type:   typeFile,
		UserId: userId,
	}
}

func NewCreateFiles(path string, name string, typeFile int32, userId string) *Files {
	return &Files{
		Path:   path,
		Name:   name,
		Type:   typeFile,
		UserId: userId,
	}
}
