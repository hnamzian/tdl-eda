package message

import (
	"encoding/json"
	"tickets/api"
	"tickets/entities"
	"tickets/message/event"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/redis/go-redis/v9"
)

func NewRouter(spreadsheetService api.SpreadsheetsClient, receiptService api.ReceiptsClient, rdb *redis.Client, watermillLogger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	h := event.NewHandler(spreadsheetService, receiptService)

	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	cancelTicketSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "cancel-ticket",
	}, watermillLogger)
	if err != nil {
		panic(err)
	}

	retry := middleware.Retry{
		MaxRetries:      10,
		InitialInterval: time.Millisecond * 100,
		MaxInterval:     time.Second,
		Multiplier:      2,
		Logger:          watermillLogger,
	}
	router.AddMiddleware(SkipMalformedEventsMiddleware)
	router.AddMiddleware(retry.Middleware)
	router.AddMiddleware(CorrelationIDMiddleware)
	router.AddMiddleware(LoggingMiddleware)

	router.AddNoPublisherHandler(
		"issue_receipt",
		"TicketBookingConfirmed",
		issueReceiptSub,
		func(msg *message.Message) error {
			if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
				return nil
			}

			var event entities.TicketBookingConfirmed
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return err
			}

			return h.IssueReceipt(msg.Context(), event)
		},
	)

	router.AddNoPublisherHandler(
		"print_ticket",
		"TicketBookingConfirmed",
		appendToTrackerSub,
		func(msg *message.Message) error {
			if msg.Metadata.Get("type") != "TicketBookingConfirmed" {
				return nil
			}
			var event entities.TicketBookingConfirmed
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return err
			}

			return h.AppendToTracker(msg.Context(), event)
		},
	)

	router.AddNoPublisherHandler(
		"cancel_ticket",
		"TicketBookingCanceled",
		cancelTicketSub,
		func(msg *message.Message) error {
			if msg.Metadata.Get("type") != "TicketBookingCanceled" {
				return nil
			}

			var event entities.TicketBookingCanceled
			err := json.Unmarshal(msg.Payload, &event)
			if err != nil {
				return err
			}

			return h.CancelTicket(msg.Context(), event)
		},
	)

	return router
}
