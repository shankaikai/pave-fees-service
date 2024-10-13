package workflow

import (
	"encore.dev/rlog"
	"go.temporal.io/sdk/workflow"
)

// BillWorkflow models the lifecycle of a bill
func BillWorkflow(ctx workflow.Context, currency string) error {
	var closed bool
	rlog.Info("Bill workflow started with currency: %s", currency)

	// Create a selector to listen for signals
	selector := workflow.NewSelector(ctx)

	// Register the signal handler for adding a line item
	closeChan := workflow.GetSignalChannel(ctx, "closeBill")
	selector.AddReceive(closeChan, func(c workflow.ReceiveChannel, more bool) {
		var signal CloseBillSignal
		c.Receive(ctx, &signal)
		rlog.Info("Received close bill signal")
		closed = true
	})

	// Workflow loop to keep listening for signals until the bill is closed
	for {
			// Wait for any of the registered events
			selector.Select(ctx)

			if closed {
					break
			}
	}

	rlog.Info("Bill workflow completed")
	return nil
}