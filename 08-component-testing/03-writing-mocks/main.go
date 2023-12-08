package main

import (
	"context"
	"sync"
	"time"
)

type IssueReceiptRequest struct {
	TicketID string `json:"ticket_id"`
	Price    Money  `json:"price"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error)
}

type ReceiptsServiceMock struct {
	lock sync.Mutex
	IssuedReceipts []IssueReceiptRequest
}

func (s *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request IssueReceiptRequest) (IssueReceiptResponse, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.IssuedReceipts = append(s.IssuedReceipts, request)
	return IssueReceiptResponse{
		ReceiptNumber: "123",
		IssuedAt:      time.Now(),
	}, nil
} 
