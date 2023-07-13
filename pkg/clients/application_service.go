package clients

import (
	"fmt"
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerClients interface {
	CreateClients(fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Clients, int, error)
	UpdateClients(id int64, fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Clients, int, error)
	DeleteClients(id int64) (int, error)
	GetClientsByID(id int64) (*Clients, int, error)
	GetAllClients() ([]*Clients, error)
	GetClientsByNit(nit string) (*Clients, int, error)
}

type service struct {
	repository ServicesClientsRepository
	txID       string
}

func NewClientsService(repository ServicesClientsRepository, TxID string) PortsServerClients {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateClients(fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Clients, int, error) {
	m := NewCreateClients(fullName, nit, banner, logoSmall, mainColor, secondColor, urlRedirect, urlApi)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Clients :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateClients(id int64, fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) (*Clients, int, error) {
	m := NewClients(id, fullName, nit, banner, logoSmall, mainColor, secondColor, urlRedirect, urlApi)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Clients :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteClients(id int64) (int, error) {
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

func (s *service) GetClientsByID(id int64) (*Clients, int, error) {
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

func (s *service) GetAllClients() ([]*Clients, error) {
	return s.repository.getAll()
}

func (s *service) GetClientsByNit(nit string) (*Clients, int, error) {
	m, err := s.repository.getByNit(nit)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByNit row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
