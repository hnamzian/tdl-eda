package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		logger.Error("Error creating subscriber", err, nil)
		panic(err)
	}

	message, err := subscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		logger.Error("Error subscribing", err, nil)
		panic(err)
	}

	for msg := range message {
		fmt.Printf("Message ID: %s - %s%\n", msg.UUID, msg.Payload)
		msg.Ack()
	}

}
