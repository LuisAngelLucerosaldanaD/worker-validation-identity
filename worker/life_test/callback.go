package life_test

import (
	"encoding/base64"
	"github.com/fatih/color"
	"log"
	"sync"
	"worker-validation-identity/infrastructure/aws_ia"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/pkg"
	"worker-validation-identity/pkg/onboarding"
)

type WorkerLifeTest struct {
	Srv *pkg.Server
}

func (w *WorkerLifeTest) StartLifeTest() {
	works, err := w.Srv.SrvOnboarding.GetAllOnboardingByStatus("life-test")
	if err != nil {
		logger.Error.Println("No se pudo obtener el listado de trabajo pendiente de prueba de vida, error: %s", err.Error())
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
		log.Printf("Worker Validation Terminated")
	}()

	close(workChan)
	wg.Wait()
}

func (w *WorkerLifeTest) doWork(work *onboarding.Onboarding) {

	user, _, err := w.Srv.SrvUser.GetUserByID(work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener los datos del usuario, error: %v", err)
		return
	}

	if user == nil {
		logger.Error.Printf("No se encontro un usuario con el id")
		return
	}

	fileSelfie, _, err := w.Srv.SrvFile.GetFileByTypeAndUserID(1, work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el registro de la selfie, error: %v", err)
		return
	}

	if fileSelfie == nil {
		logger.Error.Printf("No se pudo obtener el registro de la selfie")
		return
	}

	selfieB64, _, err := w.Srv.SrvFilesS3.GetFileByPath(fileSelfie.Path, fileSelfie.Name)
	if err != nil {
		logger.Error.Printf("No se obtener el archivo de la selfie, error: %v", err)
		return
	}

	selfieStorageBytes, err := base64.StdEncoding.DecodeString(selfieB64.Encoding)
	if err != nil {
		logger.Error.Printf("No se pudo decodificar el archivo de la selfie de base64 a bytes, error: %v", err)
		return
	}

	fileDocument, _, err := w.Srv.SrvFile.GetFileByTypeAndUserID(2, work.UserId)
	if err != nil {
		logger.Error.Printf("No se obtener el registro del lado frontal del documento, error: %v", err)
		return
	}

	if fileDocument == nil {
		logger.Error.Printf("No se obtener el registro del lado frontal del documento")
		return
	}

	documentB64, _, err := w.Srv.SrvFilesS3.GetFileByPath(fileDocument.Path, fileDocument.Name)
	if err != nil {
		logger.Error.Printf("No se pudo obtener el achivo del lado frontal del documento, error: %v", err)
		return
	}

	if documentB64 == nil {
		logger.Error.Printf("No se pudo obtener el achivo del lado frontal del documento")
		return
	}

	documentStorageBytes, err := base64.StdEncoding.DecodeString(documentB64.Encoding)
	if err != nil {
		logger.Error.Printf("No se pudo decodificar el archivo del lado frontal del documento de base64 a bytes, error: %v", err)
		return
	}

	resp, err := aws_ia.CompareFacesV2(documentStorageBytes, selfieStorageBytes)
	if err != nil {
		logger.Error.Printf("No se pudo comparar los rostros, error: %v", err)
		return
	}

	if !resp {
		_, _, err = w.Srv.SrvOnboarding.UpdateOnboarding(work.ID, work.ClientId, work.RequestId, work.UserId, "life-test-refused", work.TransactionId)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el estado de la solicitud, error: %v", err)
			return
		}

		_, _, err = w.Srv.SrvTraceability.CreateTraceability("Prueba de vida", "error",
			"La comparación entre la selfie y la foto del lado frontal del documento no coinciden o no son la misma persona", work.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo crear la trazabilidad, error: %v", err)
			return
		}
		return
	}

	_, _, err = w.Srv.SrvOnboarding.UpdateOnboarding(work.ID, work.ClientId, work.RequestId, work.UserId, "document-icr", work.TransactionId)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el estado de la solicitud, error: %v", err)
		return
	}

	_, _, err = w.Srv.SrvTraceability.CreateTraceability("Prueba de vida", "success",
		"Prueba de vida superada correctamente", work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo crear la trazabilidad, error: %v", err)
		return
	}
	return
}
