package event

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) CancelTicket(ctx context.Context, event entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Infof("Cancelling ticket")

	return h.SpreadSheetsAPI.AppendRow(
		ctx,
		"tickets-to-print",
		[]string{event.TicketID, event.CustomerEmail, event.Price.Amount, event.Price.Currency},
	)
}
