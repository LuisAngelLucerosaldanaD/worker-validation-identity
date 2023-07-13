package users

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesUsersRepository interface {
	create(m *Users) error
	update(m *Users) error
	delete(id string) error
	getByID(id string) (*Users, error)
	getAll() ([]*Users, error)
	getByEmail(email string) (*Users, error)
	getLasted(email string, limit, offset int) ([]*Users, error)
	getNotStarted() ([]*Users, error)
	getNoUploadFile(fileType int) ([]*Users, error)
	getByIdentityNumber(identityNumber int64) (*Users, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesUsersRepository {
	var s ServicesUsersRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newUsersPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
