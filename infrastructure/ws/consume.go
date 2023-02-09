package ws

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"worker-validation-identity/infrastructure/logger"
)

func ConsumeWS(jsonBytes []byte, url, method, token string) ([]byte, int, error) {

	var req http.Request

	if method == "POST" {
		resp, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
			return nil, 1, err
		}
		req = *resp
	}

	if method == "GET" {
		resp, err := http.NewRequest(method, url, nil)
		if err != nil {
			logger.Error.Printf("no se  puedo obtener respuesta: %v  -- log: ", err)
			return nil, 1, err
		}
		req = *resp
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(&req)
	if err != nil {
		logger.Error.Printf("no se  puedo enviar la petición: %v  -- log: ", err)
		return nil, resp.StatusCode, err
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
		return rsBody, resp.StatusCode, err
	}

	return rsBody, resp.StatusCode, nil
}
