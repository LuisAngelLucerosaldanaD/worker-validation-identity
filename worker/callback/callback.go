package callback

import (
	"encoding/json"
	"github.com/fatih/color"
	"log"
	"sync"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/infrastructure/ws"
	"worker-validation-identity/pkg"
	"worker-validation-identity/pkg/onboarding"
)

type WorkerCallback struct {
	Srv *pkg.Server
}

func (w *WorkerCallback) CallbackClient() {
	works, err := w.Srv.SrvOnboarding.GetOnboardingPending()
	if err != nil {
		logger.Error.Println("No se pudo obtener el listado de trabajo pendiente, error: %s", err.Error())
		return
	}

	if works == nil || len(works) < 0 {
		color.Yellow("No hay trabajo pendiente")
		return
	}

	workChan := make(chan *onboarding.Onboarding, len(works))
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
		log.Printf("Worker Terminated")
	}()

	close(workChan)
	wg.Wait()
}

func (w *WorkerCallback) doWork(work *onboarding.Onboarding) {

	client, _, err := w.Srv.SrvClient.GetClientsByID(work.ClientId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener los datos del cliente, error: %v", err)
		return
	}
	if client == nil {
		logger.Error.Printf("No se encontro un cliente con el id, error: %v", err)
		return
	}

	user, _, err := w.Srv.SrvUsers.GetUsersByID(work.UserId)
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
		_, err = w.Srv.SrvUsers.DeleteUsers(user.ID)
		if err != nil {
			logger.Error.Printf("No se pudo borrar el usuario rechazado, error: %v", err)
			return
		}
		currentStatus = "Rejected"
	}

	reqClient := RequestOnboarding{
		DocumentNumber: user.DocumentNumber,
		Status:         currentStatus,
		RequestID:      work.RequestId,
		UserID:         work.UserId,
		VerifiedAt:     work.CreatedAt,
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

	_, _, err = w.Srv.SrvOnboarding.UpdateOnboarding(work.ID, work.ClientId, work.RequestId, work.UserId, "finished")
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el registro del onboarding, error: %v", err)
		return
	}

	return
}
