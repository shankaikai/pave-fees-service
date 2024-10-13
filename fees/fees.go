package fees

import (
	"context"
	"time"

	"encore.app/fees/workflow"
	"encore.dev/beta/errs"
	"encore.dev/rlog"
	"github.com/google/uuid"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type CreateBillRequest struct {
	Currency string `json:"currency"`
}

type AddLineItemRequest struct {
	BillId 		string `json:"billId"`
	Description string `json:"description"`
	Amount      float64 `json:"amount"`
}

type AddLineItemResponse struct {
	CurrentTotal float64 `json:"currentTotal"`
	NumberOfItems int `json:"numberOfItems"`
}

type CloseBillRequest struct {
	Id 			string  `json:"id"`
}

type CloseBillResponse struct {
	Id 			string `json:"id"`
	ClosedOn 	string `json:"closedOn"`
	Bill 		workflow.Bill `json:"bill"`
}

type CreateBillResponse struct {
	Id     string `json:"id"`
}

type GetBillRequest struct {
	Id string `json:"id"`
}

type GetBillsParams struct {
	Status string `query:"status"` // open, closed
}

type GetBillsResponse struct {
	Bills []workflow.Bill `json:"bills"`
}

var SupportedCurrencies = []string{"USD", "GEL"}

// encore:api public method=POST path=/api/bill
func (s *Service) CreateBill(ctx context.Context, req *CreateBillRequest) (*CreateBillResponse, error) {
	// Validate if the currency is supported
	if !contains(SupportedCurrencies, req.Currency) {
			return nil, s.eb.Code(errs.InvalidArgument).Msg("unsupported currency, only USD or GEL").Err()
	}

	// Generate a unique ID for the bill workflow
	billWorkFlowId := uuid.New().String()

	options := client.StartWorkflowOptions{
			ID: billWorkFlowId,
			TaskQueue: billTaskQueue,
	}

	rlog.Info("Starting bill workflow", "id", billWorkFlowId)
	
	bill := workflow.Bill{
			Currency: req.Currency,
			LineItems: make([]workflow.LineItem, 0),
			TotalAmount: 0.0,
	}
	we, err := s.client.ExecuteWorkflow(ctx, options, workflow.BillWorkflow, bill)
	if err != nil {
			return nil, s.eb.Code(errs.Internal).Msg("unable to create bill").Err()
	}

	return &CreateBillResponse{
			Id:    we.GetID(),
	}, nil
}

// encore:api public method=POST path=/api/bill/close
func (s *Service) CloseBill(ctx context.Context, req *CloseBillRequest) (*CloseBillResponse, error) {
	rlog.Info("Closing bill", "id", req.Id)

	err := s.client.SignalWorkflow(ctx, req.Id, "", "closeBill", workflow.CloseBillSignal{})
	if err != nil {
			return nil, s.eb.Code(errs.Internal).Msg("unable to cancel bill workflow").Err()
	}

	// Query the workflow to get the current state
	res, err := s.client.QueryWorkflow(ctx, req.Id, "", "getBill")
	if err != nil {
			return nil, s.eb.Code(errs.Internal).Msg("unable to get bill").Err()
	}

	var bill workflow.Bill
	res.Get(&bill)

	return &CloseBillResponse{
			Id: req.Id,
			ClosedOn: time.Now().Local().Format(time.DateOnly),
			Bill: bill,
	}, nil
}

// encore:api public method=POST path=/api/bill/add
func (s *Service) AddLineItem(ctx context.Context, req *AddLineItemRequest) (*AddLineItemResponse, error) {
	if req.Amount <= 0 {
			return nil, s.eb.Code(errs.InvalidArgument).Msg("amount must be greater than 0").Err()
	}

	rlog.Info("Adding line item to bill", "description", req.Description, "amount", req.Amount)

	err := s.client.SignalWorkflow(ctx, req.BillId, "", "addLineItem", workflow.AddLineItemSignal{
			Description: req.Description,
			Amount:      req.Amount,
	})
	if err != nil {
			return nil, s.eb.Code(errs.Internal).Msg("unable to add line item to bill").Err()
	}

	bill, err := s.client.QueryWorkflow(ctx, req.BillId, "", "getBill")
	if err != nil {
			return nil, s.eb.Code(errs.Internal).Msg("unable to get bill").Err()
	}

	var b workflow.Bill
	bill.Get(&b)

	return &AddLineItemResponse{
			CurrentTotal: b.TotalAmount,
			NumberOfItems: len(b.LineItems),
	}, nil
}

// encore:api public method=GET path=/api/bill/:id
func (s *Service) GetBill(ctx context.Context, id string) (*workflow.Bill, error) {
	rlog.Info("Getting bill", "id", id)
	
	// Query the workflow to get the current state
	res, err := s.client.QueryWorkflow(ctx, id, "", "getBill")
	if err != nil {
		return nil, s.eb.Code(errs.Internal).Msg("unable to get bill").Err()
	}

	var bill workflow.Bill
	res.Get(&bill)

	return &bill, nil
}

// encore:api public method=GET path=/api/bills
func (s *Service) GetBills(ctx context.Context, params *GetBillsParams) (*GetBillsResponse, error) {
	rlog.Info("Getting bills")

	var query string
	switch params.Status {
	case "":
			query = "WorkflowType='BillWorkflow'"
	case "closed":
			query = "WorkflowType='BillWorkflow' and ExecutionStatus = 'Completed'"
	case "open":
			query = "WorkflowType='BillWorkflow' and ExecutionStatus = 'Running'"
	default:
			return nil, s.eb.Code(errs.InvalidArgument).Msg("invalid status parameter").Err()
	}
	
	options := &workflowservice.ListWorkflowExecutionsRequest{
		Query: query,
	}

	// Query the workflow to get the current state
	res, err := s.client.ListWorkflow(ctx, options)
	rlog.Info("Listing workflows", "res", res)
	if err != nil {
		rlog.Error("Error listing workflows", "error", err)
		return nil, s.eb.Code(errs.Internal).Msg("unable to get bills").Err()
	}


	var bills []workflow.Bill = make([]workflow.Bill, 0)

	for _, e := range res.Executions {
		var bill workflow.Bill
		workflowID := e.GetExecution().WorkflowId
		runID := e.GetExecution().RunId

		// Query the workflow for the bill details
		queryRes, err := s.client.QueryWorkflow(ctx, workflowID, runID, "getBill")
		if err != nil {
				rlog.Error("Error querying workflow", "workflowID", workflowID, "runID", runID, "error", err)
				continue
		}

		err = queryRes.Get(&bill)
		if err != nil {
				rlog.Error("Error getting query result", "workflowID", workflowID, "runID", runID, "error", err)
				continue
		}

		bill.Id = workflowID
		bills = append(bills, bill)
	}

	return &GetBillsResponse{Bills: bills}, nil
}

