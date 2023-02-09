package pkg

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/pkg/files"
	"worker-validation-identity/pkg/files_s3"
	"worker-validation-identity/pkg/status_request"
	"worker-validation-identity/pkg/traceability"
	"worker-validation-identity/pkg/work_validation"
)

type Server struct {
	SrvWork         work_validation.PortsServerWorkValidation
	SrvFiles        files.PortsServerFiles
	SrvFilesS3      files_s3.PortsServerFile
	SrvStatusReq    status_request.PortsServerStatusRequest
	SrvTraceability traceability.PortsServerTraceability
}

func NewServerWorker(db *sqlx.DB, txID string) *Server {
	repoWork := work_validation.FactoryStorage(db, txID)
	srvWork := work_validation.NewWorkValidationService(repoWork, txID)

	repoFiles := files.FactoryStorage(db, txID)
	srvFiles := files.NewFilesService(repoFiles, txID)

	repoS3File := files_s3.FactoryFileDocumentRepository(txID)
	srvFilesS3 := files_s3.NewFileService(repoS3File, txID)

	repoStatusReq := status_request.FactoryStorage(db, txID)
	srvStatusReq := status_request.NewStatusRequestService(repoStatusReq, txID)

	repoTraceability := traceability.FactoryStorage(db, txID)
	srvTraceability := traceability.NewTraceabilityService(repoTraceability, txID)

	return &Server{
		SrvWork:         srvWork,
		SrvFiles:        srvFiles,
		SrvFilesS3:      srvFilesS3,
		SrvStatusReq:    srvStatusReq,
		SrvTraceability: srvTraceability,
	}
}
