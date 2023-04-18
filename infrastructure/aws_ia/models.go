package aws_ia

type User struct {
	Names          string `json:"names"`
	Surname        string `json:"surname"`
	SecondSurname  string `json:"second_surname"`
	IdentityNumber string `json:"identity_number"`
}

type CompareFaceResponse struct {
	Code  int         `json:"code"`
	Data  CompareFace `json:"data"`
	Error bool        `json:"error"`
	Msg   string      `json:"msg"`
	Type  string      `json:"type"`
}

type CompareFace struct {
	DetectorBackend string  `json:"detector_backend"`
	Distance        float64 `json:"distance"`
	FacialAreas     struct {
		Img1 FacialArea `json:"img1"`
		Img2 FacialArea `json:"img2"`
	} `json:"facial_areas"`
	Model            string  `json:"model"`
	SimilarityMetric string  `json:"similarity_metric"`
	Threshold        float64 `json:"threshold"`
	Time             float64 `json:"time"`
	Verified         string  `json:"verified"`
}

type FacialArea struct {
	H int `json:"h"`
	W int `json:"w"`
	X int `json:"x"`
	Y int `json:"y"`
}
