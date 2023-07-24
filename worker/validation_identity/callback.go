package callback

import (
	"encoding/json"
	"github.com/fatih/color"
	"log"
	"sync"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/infrastructure/ws"
	"worker-validation-identity/pkg"
	"worker-validation-identity/pkg/validation_request"
)

type WorkerValidationIdentity struct {
	Srv *pkg.Server
}

func (w *WorkerValidationIdentity) SendValidationIdentity() {
	works, err := w.Srv.SrvValidationRequest.GetPendingValidationRequest()
	if err != nil {
		logger.Error.Println("No se pudo obtener el listado de trabajo de validacion de identidad, error: %s", err.Error())
		return
	}

	if works == nil || len(works) < 0 {
		color.Yellow("No hay trabajo pendiente de validacion de identidad")
		return
	}

	workChan := make(chan *validation_request.ValidationRequest, len(works))
	var wg sync.WaitGroup
	for _, work := range works {
		workChan <- work
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for workItem := range workChan {
			w.doWork(workItem)
		}
		log.Printf("Worker Validation Terminated")
	}()

	close(workChan)
	wg.Wait()
}

func (w *WorkerValidationIdentity) doWork(work *validation_request.ValidationRequest) {

	client, _, err := w.Srv.SrvClient.GetClientsByID(work.ClientId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener los datos del cliente, error: %v", err)
		return
	}
	if client == nil {
		logger.Error.Printf("No se encontro un cliente con el id, error: %v", err)
		return
	}

	user, _, err := w.Srv.SrvUsers.GetUsersByIdentityNumber(work.UserIdentification)
	if err != nil {
		logger.Error.Printf("No se pudo obtener los datos del usuario, error: %v", err)
		return
	}

	if user == nil {
		logger.Error.Printf("No se encontro un usuario con el id, error: %v", err)
		return
	}

	currentStatus := "Accepted"
	if work.Status == "refused" {
		currentStatus = "Rejected"
	}

	reqClient := RequestOnboarding{
		DocumentNumber: user.DocumentNumber,
		Status:         currentStatus,
		RequestID:      work.RequestId,
		UserID:         user.ID,
		VerifiedAt:     work.UpdatedAt,
	}

	reqBytes, _ := json.Marshal(&reqClient)
	_, code, err := ws.ConsumeWS(reqBytes, client.UrlApi, "POST", "")
	if err != nil {
		logger.Error.Printf("No se pudo enviar la petición para registra la validación de identidad, error: %v", err)
		return
	}

	if code != 200 {
		logger.Error.Printf("El servicio del cliente %s, respondio con un codigo diferente a 200, código: %d", code)
		return
	}

	// TODO pendiente guardar la respuesta del cliente en la base de datos

	_, _, err = w.Srv.SrvValidationRequest.UpdateStatusValidationRequest(work.ID, "finished")
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el registro del onboarding, error: %v", err)
		return
	}

	return
}
