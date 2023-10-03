package life_test

import (
	"fmt"

	"time"
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerLifeTest interface {
	CreateLifeTest(clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userID string, status string) (*LifeTest, int, error)
	UpdateLifeTest(id int64, clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userID string, status string) (*LifeTest, int, error)
	DeleteLifeTest(id int64) (int, error)
	GetLifeTestByID(id int64) (*LifeTest, int, error)
	GetAllLifeTest() ([]*LifeTest, error)
	GetLifeTestByClientIDAndRequestID(clientID int64, requestID string) (*LifeTest, int, error)
	UpdateStatusLifeTest(id int64, status string) (*LifeTest, int, error)
	GetAllLifeTestByUserID(userID string) ([]*LifeTest, error)
	GetLifeTestByUserID(userID string) (*LifeTest, int, error)
	GetAllLifeTestByStatus(status string) ([]*LifeTest, error)
}

type service struct {
	repository ServicesLifeTestRepository

	txID string
}

func NewLifeTestService(repository ServicesLifeTestRepository, TxID string) PortsServerLifeTest {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateLifeTest(clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userID string, status string) (*LifeTest, int, error) {
	m := NewCreateLifeTest(clientId, maxNumValidation, requestId, expiredAt, userID, status)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create LifeTest :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateLifeTest(id int64, clientId int64, maxNumValidation int, requestId string, expiredAt time.Time, userID string, status string) (*LifeTest, int, error) {
	m := NewLifeTest(id, clientId, maxNumValidation, requestId, expiredAt, userID, status)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update LifeTest :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteLifeTest(id int64) (int, error) {
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

func (s *service) GetLifeTestByID(id int64) (*LifeTest, int, error) {
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

func (s *service) GetAllLifeTest() ([]*LifeTest, error) {
	return s.repository.getAll()
}

func (s *service) GetLifeTestByClientIDAndRequestID(clientID int64, requestID string) (*LifeTest, int, error) {
	m, err := s.repository.getByClientIDAndRequestID(clientID, requestID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByClientIDAndRequestID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) UpdateStatusLifeTest(id int64, status string) (*LifeTest, int, error) {
	m := LifeTest{
		ID:     id,
		Status: status,
	}
	if err := s.repository.updateStatus(&m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update LifeTest :", err)
		return &m, 18, err
	}
	return &m, 29, nil
}

func (s *service) GetAllLifeTestByUserID(userID string) ([]*LifeTest, error) {
	return s.repository.getAllByUserId(userID)
}

func (s *service) GetLifeTestByUserID(userID string) (*LifeTest, int, error) {
	m, err := s.repository.getByUserID(userID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByUserID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllLifeTestByStatus(status string) ([]*LifeTest, error) {
	return s.repository.getAllByStatus(status)
}
