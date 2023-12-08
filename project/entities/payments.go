package entities

type PaymentRefund struct {
	TicketID       string `json:"ticket_id"`
	RefundReason   string `json:"refund_reason"`
	IdempotencyKey string `json:"idempotency_key"`
}
