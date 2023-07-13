package onboarding

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/infrastructure/logger"
)

const (
	Postgresql = "postgres"
)

type ServicesOnboardingRepository interface {
	create(m *Onboarding) error
	update(m *Onboarding) error
	delete(id string) error
	getByID(id string) (*Onboarding, error)
	getAll() ([]*Onboarding, error)
	getPending() ([]*Onboarding, error)
}

func FactoryStorage(db *sqlx.DB, txID string) ServicesOnboardingRepository {
	var s ServicesOnboardingRepository
	engine := db.DriverName()
	switch engine {
	case Postgresql:
		return newOnboardingPsqlRepository(db, txID)
	default:
		logger.Error.Println("el motor de base de datos no está implementado.", engine)
	}
	return s
}
