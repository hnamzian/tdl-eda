package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
)

type AlarmClient interface {
	StartAlarm() error
	StopAlarm() error
}

func ConsumeMessages(sub message.Subscriber, alarmClient AlarmClient) {
	messages, err := sub.Subscribe(context.Background(), "smoke_sensor")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		if msg.Payload[0] == '1' {
			err = alarmClient.StartAlarm()
		} else {
			err = alarmClient.StopAlarm()
		}
		if err != nil {
			msg.Nack()
		}
		msg.Ack()
	}

}
