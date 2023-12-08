package message

import (
	"github.com/sirupsen/logrus"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

func LoggingMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) (messages []*message.Message, err error) {
		logger := log.FromContext(msg.Context())

		logger = logger.WithField("message_uuid", msg.UUID)

		logger.Info("Handling a message")

		defer func() {
			if err != nil {
				logger.WithFields(logrus.Fields{
					"error":        err,
					"message_uuid": msg.UUID,
				}).Error("Message handling error")
			}
		}()

		return next(msg)
	}
}

func CorrelationIDMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		correlationID := msg.Metadata.Get("correlation_id")
		if correlationID == "" {
			correlationID = uuid.NewString()
		}
		ctx := log.ContextWithCorrelationID(msg.Context(), correlationID)
		ctx = log.ToContext(ctx, logrus.WithField("correlation_id", correlationID))
		msg.SetContext(ctx)

		return next(msg)
	}
}

func SkipMalformedEventsMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		if msg.UUID == "2beaf5bc-d5e4-4653-b075-2b36bbf28949" {
			return nil, nil
		}
		return next(msg)
	}
}

func SkipNoTypeEventsMiddleware(next message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		if msg.Metadata.Get("type") == "" {
			return nil, nil
		}
		return next(msg)
	}
}
