package workflow

import "time"

type AddLineItemSignal struct {
	Description string
	Amount      float64
}

type CloseBillSignal struct{}

type Bill struct {
	Id 			string `json:"id"`
	Currency   string `json:"currency"`
	LineItems []LineItem `json:"lineItems"`
	TotalAmount 	 float64 `json:"totalAmount"`
	CreatedAt *time.Time `json:"createdAt"`
	ClosedOn   *time.Time `json:"closedOn"`
}

type LineItem struct {
	Description string `json:"description"`
	Amount      float64 `json:"amount"`
	CreatedAt   *time.Time `json:"createdAt"`
}