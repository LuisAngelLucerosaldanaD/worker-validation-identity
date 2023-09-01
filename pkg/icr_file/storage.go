package icr_file

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesIcrFileRepository interface {
	create(m *IcrFile) error
	update(m *IcrFile) error
	delete(id int64) error
	getByID(id int64) (*IcrFile, error)
	getAll() ([]*IcrFile, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesIcrFileRepository {
	var s ServicesIcrFileRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newIcrFilePsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
