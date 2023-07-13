package clients

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// Clients  Model struct Clients
type Clients struct {
	ID          int64     `json:"id" db:"id" valid:"-"`
	FullName    string    `json:"full_name" db:"full_name" valid:"required"`
	Nit         string    `json:"nit" db:"nit" valid:"required"`
	Banner      string    `json:"banner" db:"banner" valid:"required"`
	LogoSmall   string    `json:"logo_small" db:"logo_small" valid:"required"`
	MainColor   string    `json:"main_color" db:"main_color" valid:"required"`
	SecondColor string    `json:"second_color" db:"second_color" valid:"required"`
	UrlRedirect string    `json:"url_redirect" db:"url_redirect" valid:"required"`
	UrlApi      string    `json:"url_api" db:"url_api" valid:"required"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewClients(id int64, fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) *Clients {
	return &Clients{
		ID:          id,
		FullName:    fullName,
		Nit:         nit,
		Banner:      banner,
		LogoSmall:   logoSmall,
		MainColor:   mainColor,
		SecondColor: secondColor,
		UrlRedirect: urlRedirect,
		UrlApi:      urlApi,
	}
}

func NewCreateClients(fullName string, nit string, banner string, logoSmall string, mainColor string, secondColor string, urlRedirect string, urlApi string) *Clients {
	return &Clients{
		FullName:    fullName,
		Nit:         nit,
		Banner:      banner,
		LogoSmall:   logoSmall,
		MainColor:   mainColor,
		SecondColor: secondColor,
		UrlRedirect: urlRedirect,
		UrlApi:      urlApi,
	}
}

func (m *Clients) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
