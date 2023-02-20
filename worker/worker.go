package worker

import (
	"encoding/base64"
	"github.com/fatih/color"
	"time"
	"worker-validation-identity/infrastructure/aws_ia"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/pkg"
)

type Worker struct {
	srv *pkg.Server
}

func NewWorker(srv *pkg.Server) IWorker {
	return &Worker{srv: srv}
}

func (w Worker) Execute() {
	e := env.NewConfiguration()
	for {
		w.CompareFace()
		time.Sleep(time.Duration(e.App.WorkerInterval) * time.Second)
	}
}

func (w Worker) CompareFace() {
	works, err := w.srv.SrvWork.GetWork()
	if err != nil {
		logger.Error.Println("No se pudo obtener el listado de trabajo pendiente, error: %s", err.Error())
		return
	}

	if len(works) < 0 {
		color.Yellow("No hay trabajo pendiente")
		return
	}

	for _, work := range works {

		selfieImg, err := w.srv.SrvFiles.GetFilesByTypeAndUserID(1, work.UserId)
		if err != nil {
			logger.Error.Println("No se pudo obtener la selfie del usuario, error: %s", err.Error())
			continue
		}

		fileSelfie, _, err := w.srv.SrvFilesS3.GetFileByPath(selfieImg.Path, selfieImg.Name)
		if err != nil {
			logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
			continue
		}

		frontDocumentImg, err := w.srv.SrvFiles.GetFilesByTypeAndUserID(2, work.UserId)
		if err != nil {
			logger.Error.Println("No se pudo obtener la imagen frontal del documento del usuario, error: %s", err.Error())
			continue
		}

		fileDocument, _, err := w.srv.SrvFilesS3.GetFileByPath(frontDocumentImg.Path, frontDocumentImg.Name)
		if err != nil {
			logger.Error.Printf("No se pudo descargar el archivo, error: %s", err.Error())
			continue
		}

		selfieBytes, err := base64.StdEncoding.DecodeString(fileSelfie.Encoding)
		if err != nil {
			logger.Error.Printf("couldn't decode selfie: %v", err)
			continue
		}

		documentFrontBytes, err := base64.StdEncoding.DecodeString(fileDocument.Encoding)
		if err != nil {
			logger.Error.Printf("couldn't decode document front: %v", err)
			continue
		}

		resp, err := aws_ia.CompareFacesV2(selfieBytes, documentFrontBytes)
		if err != nil {
			logger.Error.Printf("couldn't decode identity: %v", err)
			continue
		}

		if !resp {
			_, err = w.srv.SrvWork.UpdateWorkValidationStatus("error", work.UserId)
			if err != nil {
				logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
				continue
			}

			_, err = w.srv.SrvStatusReq.UpdateStatusRequestByUserId("rechazado", "La selfie no coincide con la persona en el documento de identidad", work.UserId)
			if err != nil {
				logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
				continue
			}

			_, _, err = w.srv.SrvTraceability.CreateTraceability("Validación de datos", "error", "La selfie no coincide con la persona en el documento de identidad", work.UserId)
			if err != nil {
				logger.Error.Printf("couldn't create traceability, error: %v", err)
				continue
			}
			continue
		}

		_, err = w.srv.SrvWork.UpdateWorkValidationStatus("ok", work.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
			continue
		}

		_, err = w.srv.SrvStatusReq.UpdateStatusRequestByUserId("pendiente", "La selfie coincide con la persona en el documento de identidad", work.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el registro, error: %s", err.Error())
			continue
		}

		_, _, err = w.srv.SrvTraceability.CreateTraceability("Validación de datos", "success", "La selfie coincide con la persona en el documento de identidad", work.UserId)
		if err != nil {
			logger.Error.Printf("couldn't create traceability, error: %v", err)
			continue
		}
		return
	}
}
