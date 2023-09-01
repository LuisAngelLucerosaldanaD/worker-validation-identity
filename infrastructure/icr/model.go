package icr

type ResponseIcrFile struct {
	Error bool      `json:"error"`
	Data  []*Letter `json:"data"`
	Code  int       `json:"code"`
	Type  int       `json:"type"`
	Msg   string    `json:"msg"`
}

type Letter struct {
	Text        string      `json:"text"`
	BoundingBox BoundingBox `json:"boundingBox"`
}
type BoundingBox struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"w"`
	H float64 `json:"h"`
}
