package event

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) AppendToTracker(ctx context.Context, event entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Infof("Appending ticket to tracker")

	if event.Price.Currency == "" {
		event.Price.Currency = "USD"
	}
	return h.SpreadSheetsAPI.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
