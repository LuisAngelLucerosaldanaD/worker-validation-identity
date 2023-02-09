package files_s3

import "worker-validation-identity/infrastructure/env"

type PortsServerFile interface {
	UploadFile(id string, originalFile string, encoding string) (*File, error)
	GetFileByPath(path, fileName string) (*ResponseFile, int, error)
}

type service struct {
	repositoryS3 ServicesFileDocumentsRepository
	txID         string
}

func NewFileService(repositoryS3 ServicesFileDocumentsRepository, TxID string) PortsServerFile {
	return &service{repositoryS3: repositoryS3, txID: TxID}
}

func (s *service) UploadFile(id string, originalFile string, encoding string) (*File, error) {
	file := NewUploadFile(id, originalFile, encoding)
	return s.repositoryS3.upload(id, file)
}

func (s *service) GetFileByPath(path, fileName string) (*ResponseFile, int, error) {
	e := env.NewConfiguration()

	rf := ResponseFile{}
	file, err := s.repositoryS3.getFile(e.Files.S3.Bucket, path, fileName)
	if err != nil {
		return nil, 0, err
	}
	rf.Encoding = file
	rf.NameDocument = fileName
	return &rf, 29, nil
}
