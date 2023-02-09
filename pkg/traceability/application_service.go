package traceability

import "worker-validation-identity/infrastructure/logger"

type PortsServerTraceability interface {
	CreateTraceability(action string, typeTrx string, description string, userId string) (*Traceability, int, error)
}

type service struct {
	repository ServicesTraceabilityRepository
	txID       string
}

func NewTraceabilityService(repository ServicesTraceabilityRepository, TxID string) PortsServerTraceability {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateTraceability(action string, typeTrx string, description string, userId string) (*Traceability, int, error) {
	m := NewCreateTraceability(action, typeTrx, description, userId)
	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Traceability :", err)
		return m, 3, err
	}
	return m, 29, nil
}
