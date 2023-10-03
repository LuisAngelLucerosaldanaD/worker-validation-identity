package file

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerFile interface {
	CreateFile(path string, name string, typeFile int32, userId string) (*File, int, error)
	UpdateFile(id int64, path string, name string, typeFile int32, userId string) (*File, int, error)
	DeleteFile(id int64) (int, error)
	GetFileByID(id int64) (*File, int, error)
	GetAllFiles() ([]*File, error)
	GetFilesByUserID(userId string) ([]*File, int, error)
	DeleteFilesByUserID(userId string) (int, error)
	GetFileByTypeAndUserID(typeFile int, userID string) (*File, int, error)
}

type service struct {
	repository ServicesFileRepository
	txID       string
}

func NewFileService(repository ServicesFileRepository, TxID string) PortsServerFile {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateFile(path string, name string, typeFile int32, userId string) (*File, int, error) {
	m := NewCreateFile(path, name, typeFile, userId)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create File :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateFile(id int64, path string, name string, typeFile int32, userId string) (*File, int, error) {
	m := NewFile(id, path, name, typeFile, userId)
	if id == 0 {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id is required"))
		return m, 15, fmt.Errorf("id is required")
	}
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update File :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteFile(id int64) (int, error) {
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

func (s *service) GetFileByID(id int64) (*File, int, error) {
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

func (s *service) GetAllFiles() ([]*File, error) {
	return s.repository.getAll()
}

func (s *service) GetFilesByUserID(userId string) ([]*File, int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("userId isn't uuid"))
		return nil, 15, fmt.Errorf("userId isn't uuid")
	}
	m, err := s.repository.getByUserID(userId)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) DeleteFilesByUserID(userId string) (int, error) {
	if !govalidator.IsUUID(userId) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("userId isn't uuid"))
		return 15, fmt.Errorf("userId isn't uuid")
	}
	if err := s.repository.deleteByUserId(userId); err != nil {
		if err.Error() == "ecatch:108" {
			return 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't delete row:", err)
		return 20, err
	}
	return 28, nil
}

func (s *service) GetFileByTypeAndUserID(typeFile int, userID string) (*File, int, error) {
	m, err := s.repository.getByTypeAndUserID(typeFile, userID)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByID row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
