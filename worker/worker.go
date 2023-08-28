package worker

import (
	"sync"
	"time"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/pkg"
	"worker-validation-identity/worker/callback"
	"worker-validation-identity/worker/life_test"
	"worker-validation-identity/worker/validation_identity"
)

type Worker struct {
	srv *pkg.Server
}

func NewWorker(srv *pkg.Server) IWorker {
	return &Worker{srv: srv}
}

func (w Worker) Execute() {
	e := env.NewConfiguration()
	callbackSrv := callback.WorkerCallback{Srv: w.srv}
	validationIdentity := validation_identity.WorkerValidationIdentity{Srv: w.srv}
	lifeTest := life_test.WorkerLifeTest{Srv: w.srv}

	var syncWorker sync.WaitGroup
	syncWorker.Add(3)
	go func() {
		defer syncWorker.Done()
		for {
			lifeTest.StartLifeTest()
			time.Sleep(time.Duration(e.App.WorkerInterval) * time.Second)
		}
	}()
	go func() {
		defer syncWorker.Done()
		for {
			callbackSrv.CallbackClient()
			time.Sleep(time.Duration(e.App.WorkerInterval) * time.Second)
		}
	}()
	go func() {
		defer syncWorker.Done()
		for {
			validationIdentity.SendValidationIdentity()
			time.Sleep(time.Duration(e.App.WorkerInterval) * time.Second)
		}
	}()
	syncWorker.Wait()
}
