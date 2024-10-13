package workflow

import "time"

type AddLineItemSignal struct {
	Description string
	Amount      float64
}

type CloseBillSignal struct{}

type Bill struct {
	Currency   string
	LineItems []LineItem
	TotalAmount 	 float64
}

type LineItem struct {
	Description string
	Amount      float64
	CreatedAt   *time.Time
}