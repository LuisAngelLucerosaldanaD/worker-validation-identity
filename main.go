package main

import (
	"github.com/fatih/color"
	"github.com/google/uuid"
	"os"
	"worker-validation-identity/infrastructure/dbx"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/pkg"
	"worker-validation-identity/worker"
)

func main() {
	c := env.NewConfiguration()
	_ = os.Setenv("AWS_ACCESS_KEY_ID", c.Aws.AWSACCESSKEYID)
	_ = os.Setenv("AWS_SECRET_ACCESS_KEY", c.Aws.AWSSECRETACCESSKEY)
	_ = os.Setenv("AWS_DEFAULT_REGION", c.Aws.AWSDEFAULTREGION)
	color.Blue("Check ID Worker v1.0.0")
	db := dbx.GetConnection()
	srv := pkg.NewServerWorker(db, uuid.New().String())
	wk := worker.NewWorker(srv)
	wk.Execute()
}
