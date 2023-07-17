package validation_request

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesValidationRequestRepository interface {
	create(m *ValidationRequest) error
	update(m *ValidationRequest) error
	delete(id int64) error
	getByID(id int64) (*ValidationRequest, error)
	getAll() ([]*ValidationRequest, error)
	getByClientIDAndRequestID(clientIid int64, requestID string) (*ValidationRequest, error)
	updateStatus(m *ValidationRequest) error
	getPending() ([]*ValidationRequest, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesValidationRequestRepository {
	var s ServicesValidationRequestRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newValidationRequestPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
