package persons

type responsePerson struct {
	Error bool    `json:"error"`
	Data  *Person `json:"data"`
	Code  int     `json:"code"`
	Type  int     `json:"type"`
	Msg   string  `json:"msg"`
}

type Person struct {
	IdentityNumber string `json:"identity_number"`
	FirstName      string `json:"first_name"`
	SecondName     string `json:"second_name"`
	Surname        string `json:"surname"`
	SecondSurname  string `json:"second_surname"`
	Particle       string `json:"particle"`
	Validity       string `json:"validity"`
	BirthDate      string `json:"birth_date"`
	ExpeditionDate string `json:"expedition_date"`
	Gender         string `json:"gender"`
}
