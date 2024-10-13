package workflow

import (
	"encore.dev/rlog"
	"go.temporal.io/sdk/workflow"
)

// BillWorkflow models the lifecycle of a bill
func BillWorkflow(ctx workflow.Context, currency string) error {
	rlog.Info("Bill workflow started with currency: %s", currency)

	// Create a selector to listen for signals
	selector := workflow.NewSelector(ctx)
	// Workflow loop to keep listening for signals until the bill is closed
	for {
			// Wait for any of the registered events
			selector.Select(ctx)
	}
}