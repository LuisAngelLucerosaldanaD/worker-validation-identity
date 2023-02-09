package status_request

import (
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerStatusRequest interface {
	UpdateStatusRequestByUserId(status string, description string, userID string) (int, error)
}

type service struct {
	repository ServicesStatusRequestRepository
	txID       string
}

func NewStatusRequestService(repository ServicesStatusRequestRepository, TxID string) PortsServerStatusRequest {
	return &service{repository: repository, txID: TxID}
}

func (s *service) UpdateStatusRequestByUserId(status string, description string, userID string) (int, error) {
	if err := s.repository.updateStatus(status, description, userID); err != nil {
		logger.Error.Println(s.txID, " - couldn't update WorkValidation :", err)
		return 18, err
	}
	return 29, nil
}
