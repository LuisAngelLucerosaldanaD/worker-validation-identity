package users

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Users  Model struct Users
type Users struct {
	ID             string     `json:"id" db:"id" valid:"required,uuid"`
	TypeDocument   *string    `json:"type_document" db:"type_document" valid:"-"`
	DocumentNumber string     `json:"document_number" db:"document_number" valid:"-"`
	ExpeditionDate *time.Time `json:"expedition_date" db:"expedition_date" valid:"-"`
	Email          string     `json:"email" db:"email" valid:"required"`
	FirstName      *string    `json:"first_name" db:"first_name" valid:"-"`
	SecondName     *string    `json:"second_name" db:"second_name" valid:"-"`
	SecondSurname  *string    `json:"second_surname" db:"second_surname" valid:"-"`
	Age            *int32     `json:"age" db:"age" valid:"-"`
	Gender         *string    `json:"gender" db:"gender" valid:"-"`
	Nationality    *string    `json:"nationality" db:"nationality" valid:"-"`
	CivilStatus    *string    `json:"civil_status" db:"civil_status" valid:"-"`
	FirstSurname   *string    `json:"first_surname" db:"first_surname" valid:"-"`
	BirthDate      *time.Time `json:"birth_date" db:"birth_date" valid:"-"`
	Country        *string    `json:"country" db:"country" valid:"-"`
	Department     *string    `json:"department" db:"department" valid:"-"`
	Cellphone      string     `json:"cellphone" db:"cellphone" valid:"-"`
	City           *string    `json:"city" db:"city" valid:"-"`
	RealIp         string     `json:"real_ip" db:"real_ip" valid:"required"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

func NewUsers(id string, typeDocument *string, documentNumber string, expeditionDate *time.Time, email string, firstName *string, secondName *string, secondSurname *string, age *int32, gender *string, nationality *string, civilStatus *string, firstSurname *string, birthDate *time.Time, country *string, department *string, city *string, realIp string, cellphone string) *Users {
	return &Users{
		ID:             id,
		TypeDocument:   typeDocument,
		DocumentNumber: documentNumber,
		ExpeditionDate: expeditionDate,
		Email:          email,
		FirstName:      firstName,
		SecondName:     secondName,
		SecondSurname:  secondSurname,
		Age:            age,
		Gender:         gender,
		Nationality:    nationality,
		CivilStatus:    civilStatus,
		FirstSurname:   firstSurname,
		BirthDate:      birthDate,
		Country:        country,
		Department:     department,
		City:           city,
		RealIp:         realIp,
		Cellphone:      cellphone,
	}
}

func (m *Users) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
