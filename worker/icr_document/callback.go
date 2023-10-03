package icr_document

import (
	"encoding/base64"
	"github.com/fatih/color"
	"log"
	"strings"
	"sync"
	"time"
	"worker-validation-identity/infrastructure/icr"
	"worker-validation-identity/infrastructure/logger"
	"worker-validation-identity/infrastructure/persons"
	"worker-validation-identity/pkg"
	"worker-validation-identity/pkg/onboarding"
	"worker-validation-identity/pkg/user"
)

type WorkerIcrDocument struct {
	Srv *pkg.Server
}

func (w *WorkerIcrDocument) StartIcrDocument() {
	works, err := w.Srv.SrvOnboarding.GetAllOnboardingByStatus("document-icr")
	if err != nil {
		logger.Error.Println("No se pudo obtener el listado de trabajo pendiente de icr, error: %s", err.Error())
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

func (w *WorkerIcrDocument) doWork(work *onboarding.Onboarding) {

	userFound, _, err := w.Srv.SrvUser.GetUserByID(work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo obtener los datos del usuario, error: %v", err)
		return
	}

	if userFound == nil {
		logger.Error.Printf("No se encontro un usuario con el id")
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

	srvIcr := icr.Icr{File: documentStorageBytes}

	letters, err := srvIcr.ProcessDocument()
	if err != nil {
		logger.Error.Printf("No se pudo realizar el proceso de icr, error: %v", err)
		return
	}

	personSrv := persons.Persons{IdentityNumber: userFound.DocumentNumber}
	userRnec, err := personSrv.GetPersonByIdentityNumber()
	if err != nil {
		logger.Error.Printf("No se pudo obtener el registro de la persona de RNec, error: %v", err)
		return
	}

	personProcess := w.ProcessIcr(letters, userRnec, userFound)

	if personProcess.Dni == "" || personProcess.SecondSurname == "" || personProcess.Name == "" || personProcess.Surname == "" {
		_, _, err = w.Srv.SrvOnboarding.UpdateOnboarding(work.ID, work.ClientId, work.RequestId, work.UserId, "icr-document-refused", work.TransactionId)
		if err != nil {
			logger.Error.Printf("No se pudo actualizar el estado de la solicitud, error: %v", err)
			return
		}

		_, _, err = w.Srv.SrvTraceability.CreateTraceability("ICR del documento", "error",
			"El proceso de extracción no se pudo realizar, se requiere validación manual", work.UserId)
		if err != nil {
			logger.Error.Printf("No se pudo crear la trazabilidad, error: %v", err)
			return
		}

		for _, letter := range letters {
			_, _, err = w.Srv.SrvIcrFile.CreateIcrFile(work.UserId, letter.Text, letter.BoundingBox.X, letter.BoundingBox.Y, letter.BoundingBox.W, letter.BoundingBox.H)
			if err != nil {
				logger.Error.Printf("No se pudo registrar el texto reconocido del icr, error: %v", err)
				continue
			}
		}
	}

	_, _, err = w.Srv.SrvOnboarding.UpdateOnboarding(work.ID, work.ClientId, work.RequestId, work.UserId, "notify-client", work.TransactionId)
	if err != nil {
		logger.Error.Printf("No se pudo actualizar el estado de la solicitud, error: %v", err)
		return
	}

	_, _, err = w.Srv.SrvTraceability.CreateTraceability("ICR del documento", "success",
		"Proceso de extracción de datos de la cedula de identidad realizada correctamente", work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo crear la trazabilidad, error: %v", err)
		return
	}

	birthDate, _ := time.Parse("02-01-2006", userRnec.BirthDate)
	age := int32(time.Now().Year() - birthDate.Year())

	nationality := "Colombia"
	_, _, err = w.Srv.SrvUser.UpdateUser(&user.User{
		ID:                 userFound.ID,
		Nickname:           userFound.Nickname,
		Email:              userFound.Email,
		Password:           userFound.Password,
		FirstName:          &userRnec.FirstName,
		SecondName:         &userRnec.SecondName,
		FirstSurname:       &userRnec.Surname,
		SecondSurname:      &userRnec.SecondSurname,
		Age:                &age,
		TypeDocument:       nil,
		DocumentNumber:     userFound.DocumentNumber,
		Cellphone:          userFound.Cellphone,
		Gender:             &userRnec.Gender,
		Nationality:        &nationality,
		RealIp:             userFound.RealIp,
		StatusId:           0,
		FailedAttempts:     0,
		LastChangePassword: userFound.LastChangePassword,
		BirthDate:          &birthDate,
		CreatedAt:          userFound.CreatedAt,
	})
	if err != nil {
		logger.Error.Printf("No se pudo actualizar los datos del usuario, error: %v", err)
		return
	}

	_, _, err = w.Srv.SrvTraceability.CreateTraceability("Actualización de datos", "info", "Actualización de los datos personales", work.UserId)
	if err != nil {
		logger.Error.Printf("No se pudo crear la trazabilidad, error: %v", err)
		return
	}

	return
}

func (w *WorkerIcrDocument) ProcessIcr(letters []*icr.Letter, userRnec *persons.Person, user *user.User) PersonICR {
	foundPerson := PersonICR{}
	for _, letter := range letters {
		if strings.ToUpper(letter.Text) == userRnec.IdentityNumber || strings.ToUpper(letter.Text) == strings.ToUpper(user.DocumentNumber) {
			foundPerson.Dni = letter.Text
			continue
		}
		if strings.ToUpper(letter.Text) == userRnec.Surname || strings.ToUpper(letter.Text) == *user.FirstSurname {
			foundPerson.Surname = letter.Text
			continue
		}
		if strings.ToUpper(letter.Text) == strings.TrimSpace(userRnec.Surname+" "+userRnec.SecondSurname) || strings.ToUpper(letter.Text) == strings.TrimSpace(*user.FirstSurname+" "+*user.SecondSurname) {
			foundPerson.Surname = letter.Text
			continue
		}
		if strings.ToUpper(letter.Text) == userRnec.SecondSurname || strings.ToUpper(letter.Text) == *user.SecondSurname {
			foundPerson.SecondSurname = letter.Text
			continue
		}
		if strings.ToUpper(letter.Text) == strings.TrimSpace(userRnec.FirstName+" "+userRnec.SecondName) || strings.ToUpper(letter.Text) == strings.TrimSpace(*user.FirstName+" "+*user.SecondName) {
			foundPerson.Name = letter.Text
			continue
		}
	}

	return foundPerson
}
