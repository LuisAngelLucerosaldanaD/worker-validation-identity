package files

import (
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerFiles interface {
	GetFilesByUserID(userId string) ([]*Files, int, error)
	GetFilesByTypeAndUserID(typeFile int, userID string) (*Files, error)
}

type service struct {
	repository ServicesFilesRepository
	txID       string
}

func NewFilesService(repository ServicesFilesRepository, TxID string) PortsServerFiles {
	return &service{repository: repository, txID: TxID}
}

func (s *service) GetFilesByUserID(userId string) ([]*Files, int, error) {
	m, err := s.repository.getByUserID(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetFilesByTypeAndUserID(typeFile int, userID string) (*Files, error) {
	m, err := s.repository.getByTypeAndUserID(typeFile, userID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, err
	}
	return m, nil
}
