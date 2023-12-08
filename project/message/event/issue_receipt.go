package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) IssueReceipt(ctx context.Context, event entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Infof("Issuing receipt")

	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}
	_, err := h.ReceiptService.IssueReceipt(ctx, entities.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price:    event.Price,
	})
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}
	return err
}
