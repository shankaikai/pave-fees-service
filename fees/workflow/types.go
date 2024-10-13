package workflow

type AddLineItemSignal struct {
	Description string
	Amount      float64
}

type CloseBillSignal struct{}

type Bill struct {
	LineItems []LineItem
}

type LineItem struct {
	Description string
	Amount      float64
	CreatedAt   string
}