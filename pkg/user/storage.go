package user

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesUserRepository interface {
	create(m *User) error
	update(m *User) error
	delete(id string) error
	getByID(id string) (*User, error)
	getAll() ([]*User, error)
	getByEmail(email string) (*User, error)
	getLasted(email string, limit, offset int) ([]*User, error)
	getNotStarted() ([]*User, error)
	getNoUploadFile(fileType int) ([]*User, error)
	getByIdentityNumber(identityNumber string) (*User, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesUserRepository {
	var s ServicesUserRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUserPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
