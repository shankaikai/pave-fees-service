package workflow

import (
	"math"
	"time"

	"encore.dev/rlog"
	"go.temporal.io/sdk/workflow"
)

const (
	CloseBill = "closeBill"
	AddLineItem = "addLineItem"
	GetBill = "getBill"
)

// BillWorkflow models the lifecycle of a bill
func BillWorkflow(ctx workflow.Context, b Bill) (Bill,error) {
	rlog.Info("Bill workflow started", "id", workflow.GetInfo(ctx).WorkflowExecution.ID, "currency", b.Currency)	

	err := workflow.SetQueryHandler(ctx, GetBill, func() (Bill, error) {
		rlog.Debug("Querying bill")
		return b, nil
	})
	if err != nil {
		rlog.Error("Error setting query handler", "error", err)
		return b, err
	}

	closed := false
	
	closeChan := workflow.GetSignalChannel(ctx, CloseBill)
	addLineItemChan := workflow.GetSignalChannel(ctx, AddLineItem)

	// Workflow loop to keep listening for signals until the bill is closed
	for {
			// Create a selector to listen for signals
			selector := workflow.NewSelector(ctx)

			// Register the signal handler for closing the bill
			selector.AddReceive(closeChan, func(c workflow.ReceiveChannel, more bool) {
				var signal CloseBillSignal
				c.Receive(ctx, &signal)
				rlog.Info("Received close bill signal")
				now := time.Now()
				b.ClosedOn = &now
				closed = true
			})

			// Register the signal handler for adding a line item
			selector.AddReceive(addLineItemChan, func(c workflow.ReceiveChannel, more bool) {
				var signal AddLineItemSignal
				c.Receive(ctx, &signal)
				rlog.Info("Received add line item signal", "description", signal.Description, "amount", signal.Amount)
				now := time.Now()
				b.AddLineItem(LineItem{
					Description: signal.Description,
					Amount:      signal.Amount,
					CreatedAt:   &now,
				})
				rlog.Info("Bill total amount updated", "totalAmount", b.TotalAmount, "lineItems", b.LineItems)
			})

			// Wait for any of the registered events
			selector.Select(ctx)

			if closed {
					break
			}
	}

	rlog.Info("Bill workflow completed", "id", workflow.GetInfo(ctx).WorkflowExecution.ID)
	return b, nil
}

func (bill *Bill) AddLineItem(item LineItem) {
	bill.LineItems = append(bill.LineItems, item)
	bill.TotalAmount += math.Ceil(item.Amount*100) / 100 // Round to 2 dp
} 