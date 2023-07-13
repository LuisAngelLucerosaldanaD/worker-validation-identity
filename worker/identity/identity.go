package identity

import (
	"encoding/base64"
	"github.com/fatih/color"
	"log"
	"sync"
	"worker-validation-identity/infrastructure/aws_ia"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/pkg"
	"worker-validation-identity/pkg/work_validation"
)

type WorkerIdentity struct {
	Srv *pkg.Server
}

func (w WorkerIdentity) CompareFace() {
	works, err := w.Srv.SrvWork.GetWork()
	if err != nil {
		logger.Error.Println("No se pudo obtener el listado de trabajo pendiente, error: %s", err.Error())
		return
	}

	if len(works) < 0 {
		color.Yellow("No hay trabajo pendiente")
		return
	}

	color.Blue("Iniciando worker de validación de identidad...")
	workChan := make(chan *work_validation.WorkValidation, len(works))
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

func (w WorkerIdentity) doWork(work *work_validation.WorkValidation) {

	selfieImg, err := w.Srv.SrvFiles.GetFilesByTypeAndUserID(1, work.UserId)
	if err != nil {
		logger.Error.Println("No se pudo obtener la selfie del usuario, error: %s", err.Error())
		return
	}

	fileSelfie, _, err := w.Srv.SrvFilesS3.GetFileByPath(selfieImg.Path, selfieImg.Name)
	if err != nil {
		logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
		return
	}

	frontDocumentImg, err := w.Srv.SrvFiles.GetFilesByTypeAndUserID(2, work.UserId)
	if err != nil {
		logger.Error.Println("No se pudo obtener la imagen frontal del documento del usuario, error: %s", err.Error())
		return
	}

	fileDocument, _, err := w.Srv.SrvFilesS3.GetFileByPath(frontDocumentImg.Path, frontDocumentImg.Name)
	if err != nil {
		logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
		return
	}

	selfieBytes, err := base64.StdEncoding.DecodeString(fileSelfie.Encoding)
	if err != nil {
		logger.Error.Printf("couldn't decode selfie: %v", err)
		return
	}

	documentFrontBytes, err := base64.StdEncoding.DecodeString(fileDocument.Encoding)
	if err != nil {
		logger.Error.Printf("couldn't decode document front: %v", err)
		return
	}

	resp, err := aws_ia.CompareFacesV2(selfieBytes, documentFrontBytes)
	if err != nil {
		logger.Error.Printf("couldn't decode identity: %v", err)
		return
	}

	if !resp {
		_, err = w.Srv.SrvWork.UpdateWorkValidationStatus("error", work.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
			return
		}

		_, err = w.Srv.SrvStatusReq.UpdateStatusRequestByUserId("rechazado", "La selfie no coincide con la persona en el documento de identidad", work.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
			return
		}

		_, _, err = w.Srv.SrvTraceability.CreateTraceability("Validación de datos", "error", "La selfie no coincide con la persona en el documento de identidad", work.UserId)
		if err != nil {
			logger.Error.Printf("couldn't create traceability, error: %v", err)
			return
		}
		return
	}

	_, err = w.Srv.SrvWork.UpdateWorkValidationStatus("ok", work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
		return
	}

	_, err = w.Srv.SrvStatusReq.UpdateStatusRequestByUserId("pendiente", "La selfie coincide con la persona en el documento de identidad", work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
		return
	}

	_, _, err = w.Srv.SrvTraceability.CreateTraceability("Validación de datos", "success", "La selfie coincide con la persona en el documento de identidad", work.UserId)
	if err != nil {
		logger.Error.Printf("couldn't create traceability, error: %v", err)
		return
	}
	return
}
