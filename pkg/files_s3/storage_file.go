package files_s3

import (
	"strings"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/infrastructure/logger"
)

const (
	S3 = "s3"
)

type ServicesFileDocumentsRepository interface {
	upload(id string, file *File) (*File, error)
	getFile(bucket, path, fileName string) (string, error)
}

func FactoryFileDocumentRepository(txID string) ServicesFileDocumentsRepository {
	var s ServicesFileDocumentsRepository
	c := env.NewConfiguration()
	repo := strings.ToLower(c.Files.Repo)
	switch repo {
	case S3:
		return newDocumentFileS3Repository(txID)
	default:
		logger.Error.Println("el repositorio de documentos no está implementado.", repo)
	}
	return s
}
