package status_request

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesStatusRequestRepository interface {
	updateStatus(status string, description string, userID string) error
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesStatusRequestRepository {
	var s ServicesStatusRequestRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newStatusRequestPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
