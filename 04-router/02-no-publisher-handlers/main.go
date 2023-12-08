package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(
		message.RouterConfig{},
		logger,
	)

	router.AddNoPublisherHandler("print-farenheit", "temperature-fahrenheit", sub, func(msg *message.Message) error {
		fmt.Printf("Temperature read: %s\n", string(msg.Payload))
		return nil
	})

	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
