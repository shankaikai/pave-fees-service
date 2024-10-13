package fees

import (
	"context"
	"fmt"

	"encore.app/fees/workflow"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var billTaskQueue = "BILL_TASK_QUEUE"

//encore:service
type Service struct {
	client client.Client
	worker worker.Worker
	eb     errs.Builder
}

func initService() (*Service, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("unable to create temporal client: %v", err)
	}

	rlog.Info("Connected to Temporal server")

	w := worker.New(c, billTaskQueue, worker.Options{})

	w.RegisterWorkflow(workflow.BillWorkflow)

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("unable to start worker: %v", err)
	}

	rlog.Info("Started worker for bill workflow")

	return &Service{client: c, worker: w, eb: *errs.B()}, nil
}

func (s *Service) Shutdown(force context.Context) {
	s.client.Close()
	s.worker.Stop()
}