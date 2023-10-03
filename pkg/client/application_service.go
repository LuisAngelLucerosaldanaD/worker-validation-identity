package client

import (
	"fmt"

	"worker-validation-identity/infrastructure/logger"
)

type PortsServerClient interface {
	CreateClient(fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Client, int, error)
	UpdateClient(id int64, fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Client, int, error)
	DeleteClient(id int64) (int, error)
	GetClientByID(id int64) (*Client, int, error)
	GetAllClients() ([]*Client, error)
	GetClientByNit(nit string) (*Client, int, error)
}

type service struct {
	repository ServicesClientRepository

	txID string
}

func NewClientService(repository ServicesClientRepository, TxID string) PortsServerClient {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateClient(fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Client, int, error) {
	m := NewCreateClient(fullName, nit, banner, logoSmall, mainColor, secondColor, urlRedirect, urlApi)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Client :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateClient(id int64, fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Client, int, error) {
	m := NewClient(id, fullName, nit, banner, logoSmall, mainColor, secondColor, urlRedirect, urlApi)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Client :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteClient(id int64) (int, error) {
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

func (s *service) GetClientByID(id int64) (*Client, int, error) {
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

func (s *service) GetAllClients() ([]*Client, error) {
	return s.repository.getAll()
}

func (s *service) GetClientByNit(nit string) (*Client, int, error) {
	m, err := s.repository.getByNit(nit)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByNit row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
