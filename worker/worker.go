package worker

import (
	"sync"
	"time"
	"worker-validation-identity/infrastructure/env"
	"worker-validation-identity/pkg"
	"worker-validation-identity/worker/callback"
	"worker-validation-identity/worker/identity"
)

type Worker struct {
	srv *pkg.Server
}

func NewWorker(srv *pkg.Server) IWorker {
	return &Worker{srv: srv}
}

func (w Worker) Execute() {
	e := env.NewConfiguration()
	callbackSrv := callback.WorkerCallback{
		Srv: w.srv,
	}
	identitySrv := identity.WorkerIdentity{
		Srv: w.srv,
	}
	for {
		var syncWorker sync.WaitGroup
		syncWorker.Add(2)
		go func() {
			defer syncWorker.Done()
			identitySrv.CompareFace()
		}()
		go func() {
			defer syncWorker.Done()
			callbackSrv.CallbackClient()
		}()
		syncWorker.Wait()
		time.Sleep(time.Duration(e.App.WorkerInterval) * time.Second)
	}
}
