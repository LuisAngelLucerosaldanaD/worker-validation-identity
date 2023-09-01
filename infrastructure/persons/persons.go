package persons

import (
	"encoding/json"
	"fmt"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/infrastructure/ws"
)

type Persons struct {
	IdentityNumber string
}

func (p *Persons) GetPersonByIdentityNumber() (*Person, error) {

	resPerson := responsePerson{}
	e := env.NewConfiguration()
	resWs, code, err := ws.ConsumeWS(nil, e.RNecApi.Dni+p.IdentityNumber, "GET", "")
	if err != nil || code != 200 {
		logger.Error.Printf("No se pudo obtener la persona por el número de identificación: %v", err)
		return nil, err
	}

	err = json.Unmarshal(resWs, &resPerson)
	if err != nil {
		logger.Error.Printf("No se pudo parsear la respuesta: %v", err)
		return nil, err
	}

	if resPerson.Error {
		logger.Error.Printf("No se pudo obtener los datos de la persona, mensaje: ", resPerson.Msg)
		return nil, fmt.Errorf(resPerson.Msg)
	}

	if resPerson.Data == nil {
		logger.Error.Printf("No se encontro a la persona por su número de identificacion, mensaje: ", resPerson.Msg)
		return nil, fmt.Errorf("no se encontro a la persona por su número de identificacion")
	}

	return resPerson.Data, nil
}
