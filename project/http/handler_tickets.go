package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatus `json:"tickets"`
}

type TicketStatus struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
}

type Tickethandler struct {
	publisher message.Publisher
}

func NewTicketHandler(publisher message.Publisher) *Tickethandler {
	return &Tickethandler{
		publisher: publisher,
	}
}

func (h *Tickethandler) PostTicketStatus(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			msg.Metadata.Set("type", "TicketBookingConfirmed")

			err = h.publisher.Publish("TicketBookingConfirmed", msg)
			if err != nil {
				return err
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewEventHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := message.NewMessage(watermill.NewUUID(), payload)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))
			msg.Metadata.Set("type", "TicketBookingConfirmed")

			err = h.publisher.Publish("TicketBookingCanceled", msg)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}

