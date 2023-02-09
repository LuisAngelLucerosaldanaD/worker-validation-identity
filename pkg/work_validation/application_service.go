package work_validation

import (
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerWorkValidation interface {
	GetWork() ([]*WorkValidation, error)
	UpdateWorkValidationStatus(status string, userID string) (int, error)
}

type service struct {
	repository ServicesWorkValidationRepository
	txID       string
}

func NewWorkValidationService(repository ServicesWorkValidationRepository, TxID string) PortsServerWorkValidation {
	return &service{repository: repository, txID: TxID}
}

func (s *service) GetWork() ([]*WorkValidation, error) {
	m, err := s.repository.getPending()
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserId row:", err)
		return nil, err
	}
	return m, nil
}

func (s *service) UpdateWorkValidationStatus(status string, userID string) (int, error) {
	if err := s.repository.updateStatus(status, userID); err != nil {
		logger.Error.Println(s.txID, " - couldn't update WorkValidation :", err)
		return 18, err
	}
	return 29, nil
}
