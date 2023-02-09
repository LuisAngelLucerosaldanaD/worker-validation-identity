package traceability

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesTraceabilityRepository interface {
	create(m *Traceability) error
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesTraceabilityRepository {
	var s ServicesTraceabilityRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newTraceabilityPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
