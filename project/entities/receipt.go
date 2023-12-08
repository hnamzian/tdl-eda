package entities

import "time"

type VoidReceipt struct {
	TicketID       string `json:"ticket_id"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

type IssueReceiptRequest struct {
	TicketID string
	Price    Money
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"receipt_number"`
	IssuedAt      time.Time `json:"issued_at"`
}
