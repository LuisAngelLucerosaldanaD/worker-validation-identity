package work_validation

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesWorkValidationRepository interface {
	getPending() ([]*WorkValidation, error)
	updateStatus(status string, userID string) error
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesWorkValidationRepository {
	var s ServicesWorkValidationRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newWorkValidationPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
