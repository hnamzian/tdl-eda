package event

import (
	"context"
	"tickets/entities"
)

type Handler struct {
	SpreadSheetsAPI SpreadSheetsAPI
	ReceiptService  ReciptService
}

func NewHandler(spreadSheetsAPI SpreadSheetsAPI, receiptService ReciptService) *Handler {
	if spreadSheetsAPI == nil {
		panic("spreadsheetService is nil")
	}
	if receiptService == nil {
		panic("receiptService is nil")
	}
	return &Handler{
		SpreadSheetsAPI: spreadSheetsAPI,
		ReceiptService:  receiptService,
	}
}

type SpreadSheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReciptService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}
