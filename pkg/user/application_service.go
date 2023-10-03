package user

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"worker-validation-identity/infrastructure/logger"
)

type PortsServerUser interface {
	CreateUser(user *User) (*User, int, error)
	UpdateUser(user *User) (*User, int, error)
	DeleteUser(id string) (int, error)
	GetUserByID(id string) (*User, int, error)
	GetAllUser() ([]*User, error)
	GetUserByEmail(email string) (*User, int, error)
	GetAllUserLasted(email string, limit, offset int) ([]*User, error)
	GetAllNotStarted() ([]*User, error)
	GetAllNotUploadFile(fileType int) ([]*User, error)
	GetUserByIdentityNumber(identityNumber string) (*User, int, error)
}

type service struct {
	repository ServicesUserRepository
	txID       string
}

func NewUsersService(repository ServicesUserRepository, TxID string) PortsServerUser {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateUser(user *User) (*User, int, error) {
	if valid, err := user.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return user, 15, err
	}

	if err := s.repository.create(user); err != nil {
		if err.Error() == "ecatch:108" {
			return user, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create User :", err)
		return user, 3, err
	}
	return user, 29, nil
}

func (s *service) UpdateUser(user *User) (*User, int, error) {
	if valid, err := user.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return user, 15, err
	}
	if err := s.repository.update(user); err != nil {
		logger.Error.Println(s.txID, " - couldn't update User :", err)
		return user, 18, err
	}
	return user, 29, nil
}

func (s *service) DeleteUser(id string) (int, error) {
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

func (s *service) GetUserByID(id string) (*User, int, error) {
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

func (s *service) GetAllUser() ([]*User, error) {
	return s.repository.getAll()
}

func (s *service) GetUserByEmail(email string) (*User, int, error) {
	if !govalidator.IsEmail(email) {
		logger.Error.Println(s.txID, " - don't meet validations:", fmt.Errorf("id isn't email valid"))
		return nil, 15, fmt.Errorf("id isn't email valid")
	}
	m, err := s.repository.getByEmail(email)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByEmail row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}

func (s *service) GetAllUserLasted(email string, limit, offset int) ([]*User, error) {
	return s.repository.getLasted(email, limit, offset)
}

func (s *service) GetAllNotStarted() ([]*User, error) {
	return s.repository.getNotStarted()
}

func (s *service) GetAllNotUploadFile(fileType int) ([]*User, error) {
	return s.repository.getNoUploadFile(fileType)
}

func (s *service) GetUserByIdentityNumber(identityNumber string) (*User, int, error) {
	m, err := s.repository.getByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByIdentityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
