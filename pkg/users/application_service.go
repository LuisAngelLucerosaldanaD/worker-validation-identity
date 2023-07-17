package users

import (
	"fmt"
	"time"
	"worker-validation-identity/infrastructure/logger"

	"github.com/asaskevich/govalidator"
)

type PortsServerUsers interface {
	CreateUsers(id string, typeDocument *string, documentNumber string, expeditionDate *time.Time, email string, firstName *string, secondName *string, secondSurname *string, age *int32, gender *string, nationality *string, civilStatus *string, firstSurname *string, birthDate *time.Time, country *string, department *string, city *string, realIp string, cellphone string) (*Users, int, error)
	UpdateUsers(id string, typeDocument *string, documentNumber string, expeditionDate *time.Time, email string, firstName *string, secondName *string, secondSurname *string, age *int32, gender *string, nationality *string, civilStatus *string, firstSurname *string, birthDate *time.Time, country *string, department *string, city *string, realIp string, cellphone string) (*Users, int, error)
	DeleteUsers(id string) (int, error)
	GetUsersByID(id string) (*Users, int, error)
	GetAllUsers() ([]*Users, error)
	GetUsersByEmail(email string) (*Users, int, error)
	GetAllUsersLasted(email string, limit, offset int) ([]*Users, error)
	GetAllNotStarted() ([]*Users, error)
	GetAllNotUploadFile(fileType int) ([]*Users, error)
	GetUsersByIdentityNumber(identityNumber string) (*Users, int, error)
}

type service struct {
	repository ServicesUsersRepository
	txID       string
}

func NewUsersService(repository ServicesUsersRepository, TxID string) PortsServerUsers {
	return &service{repository: repository, txID: TxID}
}

func (s *service) CreateUsers(id string, typeDocument *string, documentNumber string, expeditionDate *time.Time, email string, firstName *string, secondName *string, secondSurname *string, age *int32, gender *string, nationality *string, civilStatus *string, firstSurname *string, birthDate *time.Time, country *string, department *string, city *string, realIp string, cellphone string) (*Users, int, error) {
	m := NewUsers(id, typeDocument, documentNumber, expeditionDate, email, firstName, secondName, secondSurname, age, gender, nationality, civilStatus, firstSurname, birthDate, country, department, city, realIp, cellphone)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}

	if err := s.repository.create(m); err != nil {
		if err.Error() == "ecatch:108" {
			return m, 108, nil
		}
		logger.Error.Println(s.txID, " - couldn't create Users :", err)
		return m, 3, err
	}
	return m, 29, nil
}

func (s *service) UpdateUsers(id string, typeDocument *string, documentNumber string, expeditionDate *time.Time, email string, firstName *string, secondName *string, secondSurname *string, age *int32, gender *string, nationality *string, civilStatus *string, firstSurname *string, birthDate *time.Time, country *string, department *string, city *string, realIp string, cellphone string) (*Users, int, error) {
	m := NewUsers(id, typeDocument, documentNumber, expeditionDate, email, firstName, secondName, secondSurname, age, gender, nationality, civilStatus, firstSurname, birthDate, country, department, city, realIp, cellphone)
	if valid, err := m.valid(); !valid {
		logger.Error.Println(s.txID, " - don't meet validations:", err)
		return m, 15, err
	}
	if err := s.repository.update(m); err != nil {
		logger.Error.Println(s.txID, " - couldn't update Users :", err)
		return m, 18, err
	}
	return m, 29, nil
}

func (s *service) DeleteUsers(id string) (int, error) {
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

func (s *service) GetUsersByID(id string) (*Users, int, error) {
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

func (s *service) GetAllUsers() ([]*Users, error) {
	return s.repository.getAll()
}

func (s *service) GetUsersByEmail(email string) (*Users, int, error) {
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

func (s *service) GetAllUsersLasted(email string, limit, offset int) ([]*Users, error) {
	return s.repository.getLasted(email, limit, offset)
}

func (s *service) GetAllNotStarted() ([]*Users, error) {
	return s.repository.getNotStarted()
}

func (s *service) GetAllNotUploadFile(fileType int) ([]*Users, error) {
	return s.repository.getNoUploadFile(fileType)
}

func (s *service) GetUsersByIdentityNumber(identityNumber string) (*Users, int, error) {
	m, err := s.repository.getByIdentityNumber(identityNumber)
	if err != nil {
		logger.Error.Println(s.txID, " - couldn`t getByIdentityNumber row:", err)
		return nil, 22, err
	}
	return m, 29, nil
}
