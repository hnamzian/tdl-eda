package service

import (
	"context"
	"net/http"
	"os"
	"tickets/api"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.Init(logrus.InfoLevel)
}

type Service struct {
	watermillRouter *watermillMessage.Router
	echoRouter      *echo.Echo
}

func New(
	redisClient *redis.Client,
	spreadsheetsService event.SpreadSheetsAPI,
	receiptsService event.ReciptService,
) Service {
	log.Init(logrus.InfoLevel)

	clients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	receiptsClient := api.NewReceiptsClient(clients)
	spreadsheetsClient := api.NewSpreadsheetsClient(clients)

	watermillLogger := log.NewWatermill(logrus.NewEntry(logrus.StandardLogger()))

	rdb := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	publisher := message.NewRedisPublisher(rdb, watermillLogger)
	router := message.NewRouter(spreadsheetsClient, receiptsClient, rdb, watermillLogger)

	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	return Service{
		watermillRouter: router,
		echoRouter:      echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
