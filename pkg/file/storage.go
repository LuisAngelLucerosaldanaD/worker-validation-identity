package file

import (
	"github.com/jmoiron/sqlx"

	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesFileRepository interface {
	create(m *File) error
	update(m *File) error
	delete(id int64) error
	getByID(id int64) (*File, error)
	getAll() ([]*File, error)
	getByUserID(userID string) ([]*File, error)
	deleteByUserId(userId string) error
	getByTypeAndUserID(typeFile int, userID string) (*File, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesFileRepository {
	var s ServicesFileRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newFilePsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
