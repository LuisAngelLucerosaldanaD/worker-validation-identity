package files

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesFilesRepository interface {
	getByUserID(userID string) ([]*Files, error)
	getByTypeAndUserID(typeFile int, userID string) (*Files, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesFilesRepository {
	var s ServicesFilesRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newFilesPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
