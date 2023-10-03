package pkg

import (
	"github.com/jmoiron/sqlx"
	"worker-validation-identity/pkg/client"
	"worker-validation-identity/pkg/file"
	"worker-validation-identity/pkg/files_s3"
	"worker-validation-identity/pkg/icr_file"
	"worker-validation-identity/pkg/life_test"
	"worker-validation-identity/pkg/onboarding"
	"worker-validation-identity/pkg/traceability"
	"worker-validation-identity/pkg/user"
)

type Server struct {
	SrvFile         file.PortsServerFile
	SrvFilesS3      files_s3.PortsServerFile
	SrvTraceability traceability.PortsServerTraceability
	SrvOnboarding   onboarding.PortsServerOnboarding
	SrvClient       client.PortsServerClient
	SrvIcrFile      icr_file.PortsServerIcrFile
	SrvLifeTest     life_test.PortsServerLifeTest
	SrvUser         user.PortsServerUser
}

func NewServerWorker(db *sqlx.DB, txID string) *Server {
	repoFile := file.FactoryStorage(db, txID)
	srvFile := file.NewFileService(repoFile, txID)

	repoS3File := files_s3.FactoryFileDocumentRepository(txID)
	srvFilesS3 := files_s3.NewFileService(repoS3File, txID)

	repoTraceability := traceability.FactoryStorage(db, txID)
	srvTraceability := traceability.NewTraceabilityService(repoTraceability, txID)

	repoOnboarding := onboarding.FactoryStorage(db, txID)
	srvOnboarding := onboarding.NewOnboardingService(repoOnboarding, txID)

	repoClient := client.FactoryStorage(db, txID)
	srvClient := client.NewClientService(repoClient, txID)

	repoIcrFile := icr_file.FactoryStorage(db, txID)
	srvIcrFile := icr_file.NewIcrFileService(repoIcrFile, txID)

	repoLifeTest := life_test.FactoryStorage(db, txID)
	srvLifeTest := life_test.NewLifeTestService(repoLifeTest, txID)

	repoUser := user.FactoryStorage(db, txID)
	srvUser := user.NewUsersService(repoUser, txID)

	return &Server{
		SrvFile:         srvFile,
		SrvFilesS3:      srvFilesS3,
		SrvTraceability: srvTraceability,
		SrvOnboarding:   srvOnboarding,
		SrvClient:       srvClient,
		SrvIcrFile:      srvIcrFile,
		SrvLifeTest:     srvLifeTest,
		SrvUser:         srvUser,
	}
}
