package user

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// User  Model struct User
type User struct {
	ID                 string     `json:"id" db:"id" valid:"required,uuid"`
	Nickname           string     `json:"nickname" db:"nickname" valid:"required"`
	Email              string     `json:"email" db:"email" valid:"required"`
	Password           string     `json:"password" db:"password" valid:"required"`
	FirstName          *string    `json:"first_name" db:"first_name" valid:"-"`
	SecondName         *string    `json:"second_name" db:"second_name" valid:"-"`
	FirstSurname       *string    `json:"first_surname" db:"first_surname" valid:"-"`
	SecondSurname      *string    `json:"second_surname" db:"second_surname" valid:"-"`
	Age                *int32     `json:"age" db:"age" valid:"-"`
	TypeDocument       *string    `json:"type_document" db:"type_document"`
	DocumentNumber     string     `json:"document_number" db:"document_number" valid:"required"`
	Cellphone          string     `json:"cellphone" db:"cellphone"`
	Gender             *string    `json:"gender" db:"gender"`
	Nationality        *string    `json:"nationality" db:"nationality"`
	Country            *string    `json:"country" db:"country"`
	Department         *string    `json:"department" db:"department"`
	City               *string    `json:"city" db:"city"`
	RealIp             string     `json:"real_ip" db:"real_ip"`
	StatusId           int32      `json:"status_id" db:"status_id"`
	FailedAttempts     int32      `json:"failed_attempts" db:"failed_attempts"`
	BlockDate          *time.Time `json:"block_date" db:"block_date"`
	DisabledDate       *time.Time `json:"disabled_date" db:"disabled_date"`
	LastLogin          *time.Time `json:"last_login" db:"last_login"`
	LastChangePassword *time.Time `json:"last_change_password" db:"last_change_password"`
	BirthDate          *time.Time `json:"birth_date" db:"birth_date"`
	VerifiedCode       *string    `json:"verified_code" db:"verified_code"`
	IsDeleted          bool       `json:"is_deleted" db:"is_deleted"`
	DeletedAt          *time.Time `json:"deleted_at" db:"deleted_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

func NewUser(id string, nickname string, email string, password string, firstName *string, secondName *string,
	firstSurname *string, secondSurname *string, age *int32, typeDocument *string, documentNumber string,
	cellphone string, gender *string, nationality *string, country *string, department *string, city *string,
	realIp string, statusId int32, failedAttempts int32, blockDate *time.Time, disabledDate *time.Time,
	lastLogin *time.Time, lastChangePassword *time.Time, birthDate *time.Time, verifiedCode *string, isDeleted bool,
	deletedAt *time.Time) *User {
	return &User{
		ID:                 id,
		Nickname:           nickname,
		Email:              email,
		Password:           password,
		FirstName:          firstName,
		SecondName:         secondName,
		FirstSurname:       firstSurname,
		SecondSurname:      secondSurname,
		Age:                age,
		TypeDocument:       typeDocument,
		DocumentNumber:     documentNumber,
		Cellphone:          cellphone,
		Gender:             gender,
		Nationality:        nationality,
		Country:            country,
		Department:         department,
		City:               city,
		RealIp:             realIp,
		StatusId:           statusId,
		FailedAttempts:     failedAttempts,
		BlockDate:          blockDate,
		DisabledDate:       disabledDate,
		LastLogin:          lastLogin,
		LastChangePassword: lastChangePassword,
		BirthDate:          birthDate,
		VerifiedCode:       verifiedCode,
		IsDeleted:          isDeleted,
		DeletedAt:          deletedAt,
	}
}

func (m *User) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
