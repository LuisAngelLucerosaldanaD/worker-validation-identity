package icr_file

import (
	"fmt"
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerIcrFile interface {
	CreateIcrFile(userId string, text string, x float64, y float64, w float64, h float64) (*IcrFile, int, error)
	UpdateIcrFile(id int64, userId string, text string, x float64, y float64, w float64, h float64) (*IcrFile, int, error)
	DeleteIcrFile(id int64) (int, error)
	GetIcrFileByID(id int64) (*IcrFile, int, error)
	GetAllIcrFile() ([]*IcrFile, error)
}

type service struct {
	repository ServicesIcrFileRepository
	txID       string
}

func NewIcrFileService(repository ServicesIcrFileRepository, TxID string) PortsServerIcrFile {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateIcrFile(userId string, text string, x float64, y float64, w float64, h float64) (*IcrFile, int, error) {
	m := NewCreateIcrFile(userId, text, x, y, w, h)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create IcrFile :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateIcrFile(id int64, userId string, text string, x float64, y float64, w float64, h float64) (*IcrFile, int, error) {
	m := NewIcrFile(id, userId, text, x, y, w, h)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update IcrFile :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteIcrFile(id int64) (int, error) {
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

func (s *service) GetIcrFileByID(id int64) (*IcrFile, int, error) {
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

func (s *service) GetAllIcrFile() ([]*IcrFile, error) {
	return s.repository.getAll()
}
