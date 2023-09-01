package icr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"worker-validation-identity/infrastructure/aws_ia"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/infrastructure/logger"
)

type Icr struct {
	File []byte
}

func (i *Icr) ProcessDocument() ([]*Letter, error) {
	e := env.NewConfiguration()
	documentReader := bytes.NewReader(i.File)
	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)

	face1Writer, err := multipartWriter.CreateFormFile("img-encoding", "dni"+aws_ia.GetExtensionFromBytes(i.File))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(face1Writer, documentReader); err != nil {
		return nil, err
	}

	multipartWriter.Close()

	request, err := http.NewRequest(http.MethodPost, e.FaceApi.CompareFace, body)
	request.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	httpClient := &http.Client{}

	resp, err := httpClient.Do(request)
	if err != nil {
		logger.Error.Printf("no se  puedo enviar la petición: %v  -- log: ", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error.Printf("no se pudo ejecutar defer body close: %v  -- log: ", err)
		}
	}(resp.Body)

	rsBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
		return nil, err
	}

	resIcr := ResponseIcrFile{}

	err = json.Unmarshal(rsBody, &resIcr)
	if err != nil {
		logger.Error.Printf("no se pudo parsear la respuesta del servio de comparacion de rostros: %v  -- log: ", err)
		return nil, err
	}

	if resIcr.Error {
		logger.Error.Printf("error al consumir el servicio de icr: " + resIcr.Msg)
		return nil, fmt.Errorf(resIcr.Msg)
	}

	if resIcr.Data == nil || len(resIcr.Data) < 0 {
		return nil, fmt.Errorf("no se pudo extraer ningun caracter del archivo procesado, verifique que sea una imagen valida")
	}

	return resIcr.Data, nil
}
