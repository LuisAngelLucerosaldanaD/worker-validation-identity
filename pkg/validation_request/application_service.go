package validation_request

import (
	"fmt"
	"worker-validation-identity/infrastructure/logger"

	"time"
)

type PortsServerValidationRequest interface {
	CreateValidationRequest(clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) (*ValidationRequest, int, error)
	UpdateValidationRequest(id int64, clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) (*ValidationRequest, int, error)
	DeleteValidationRequest(id int64) (int, error)
	GetValidationRequestByID(id int64) (*ValidationRequest, int, error)
	GetAllValidationRequest() ([]*ValidationRequest, error)
	GetValidationRequestByClientIDAndRequestID(clientID int64, requestID string) (*ValidationRequest, int, error)
	UpdateStatusValidationRequest(id int64, status string) (*ValidationRequest, int, error)
	GetPendingValidationRequest() ([]*ValidationRequest, error)
}

type service struct {
	repository ServicesValidationRequestRepository
	txID       string
}

func NewValidationRequestService(repository ServicesValidationRequestRepository, TxID string) PortsServerValidationRequest {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateValidationRequest(clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) (*ValidationRequest, int, error) {
	m := NewCreateValidationRequest(clientId, maxNumValidation, requestId, expiredAt, userIdentification, status)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create ValidationRequest :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateValidationRequest(id int64, clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userIdentification string, status string) (*ValidationRequest, int, error) {
	m := NewValidationRequest(id, clientId, maxNumValidation, requestId, expiredAt, userIdentification, status)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update ValidationRequest :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteValidationRequest(id int64) (int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return 15, fmt.Errorf("id is required")
	}

	if err := s.repository.delete(id); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't update row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetValidationRequestByID(id int64) (*ValidationRequest, int, error) {
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return nil, 15, fmt.Errorf("id is required")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllValidationRequest() ([]*ValidationRequest, error) {
	return s.repository.getAll()
}

func (s *service) GetValidationRequestByClientIDAndRequestID(clientID int64, requestID string) (*ValidationRequest, int, error) {
	m, err := s.repository.getByClientIDAndRequestID(clientID, requestID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByClientIDAndRequestID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) UpdateStatusValidationRequest(id int64, status string) (*ValidationRequest, int, error) {
	m := ValidationRequest{
		ID:     id,
		Status: status,
	}
	if err := s.repository.updateStatus(&m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update ValidationRequest :", err)
		return &m, 18, err
	}
	return &m, 29, nil
}

func (s *service) GetPendingValidationRequest() ([]*ValidationRequest, error) {
	return s.repository.getPending()
}
