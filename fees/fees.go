package fees

import (
	"context"

	"encore.app/fees/workflow"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
)

type CreateBillRequest struct {
	Currency string
}

type AddLineItemRequest struct {
	Description string
	Amount      float64
}

type CloseBillRequest struct {
	Id 			string
}

type CreateBillResponse struct {
	Id     string
}

// encore:api public method=POST path=/api/bill
func (s *Service) CreateBill(ctx context.Context, req *CreateBillRequest) (*CreateBillResponse, error) {
	// Generate a unique ID for the bill workflow
	billWorkFlowId := uuid.New().String()

	options := client.StartWorkflowOptions{
			ID: billWorkFlowId,
			TaskQueue: billTaskQueue,
	}

	rlog.Info("Starting bill workflow", "id", billWorkFlowId)

	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.BillWorkflow, req.Currency)
	if err != nil {
			return nil, s.eb.Code(errs.Internal).Msg("unable to start bill workflow").Err()
	}

	return &CreateBillResponse{
			Id:    we.GetID(),
	}, nil
}