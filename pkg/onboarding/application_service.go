package onboarding

import (
	"fmt"
	"worker-validation-identity/infrastructure/logger"

	"github.com/asaskevich/govalidator"
)

type PortsServerOnboarding interface {
	CreateOnboarding(id string, clientId int64, requestId string, userId string, status string) (*Onboarding, int, error)
	UpdateOnboarding(id string, clientId int64, requestId string, userId string, status string) (*Onboarding, int, error)
	DeleteOnboarding(id string) (int, error)
	GetOnboardingByID(id string) (*Onboarding, int, error)
	GetAllOnboarding() ([]*Onboarding, error)
	GetOnboardingPending() ([]*Onboarding, error)
}

type service struct {
	repository ServicesOnboardingRepository
	txID       string
}

func NewOnboardingService(repository ServicesOnboardingRepository, TxID string) PortsServerOnboarding {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateOnboarding(id string, clientId int64, requestId string, userId string, status string) (*Onboarding, int, error) {
	m := NewOnboarding(id, clientId, requestId, userId, status)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Onboarding :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateOnboarding(id string, clientId int64, requestId string, userId string, status string) (*Onboarding, int, error) {
	m := NewOnboarding(id, clientId, requestId, userId, status)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Onboarding :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteOnboarding(id string) (int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return 15, fmt.Errorf("id isn't uuid")
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

func (s *service) GetOnboardingByID(id string) (*Onboarding, int, error) {
	if !govalidator.IsUUID(id) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't uuid"))
		return nil, 15, fmt.Errorf("id isn't uuid")
	}
	m, err := s.repository.getByID(id)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllOnboarding() ([]*Onboarding, error) {
	return s.repository.getAll()
}

func (s *service) GetOnboardingPending() ([]*Onboarding, error) {
	return s.repository.getPending()
}
