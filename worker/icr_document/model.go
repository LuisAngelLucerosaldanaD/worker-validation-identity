package icr_document

type PersonICR struct {
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	SecondSurname string `json:"second_surname"`
	Dni           string `json:"dni"`
}
