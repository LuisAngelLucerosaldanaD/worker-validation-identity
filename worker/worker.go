package worker

import (
	"sync"
	"time"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/pkg"
	"worker-validation-identity/worker/callback"
	"worker-validation-identity/worker/identity"
	callback2 "worker-validation-identity/worker/validation_identity"
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
	identitySrv := identity.WorkerIdentity{Srv: w.srv}
	validationIdentity := callback2.WorkerValidationIdentity{Srv: w.srv}

	var syncWorker sync.WaitGroup
	syncWorker.Add(3)
	go func() {
		defer syncWorker.Done()
		for {
			identitySrv.CompareFace()
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
