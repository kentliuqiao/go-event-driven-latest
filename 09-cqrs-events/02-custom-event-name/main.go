package main

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewEventBus(pub message.Publisher) (*cqrs.EventBus, error) {
	return cqrs.NewEventBusWithConfig(pub, cqrs.EventBusConfig{
		Marshaler: cqrs.JSONMarshaler{
			GenerateName: cqrs.StructName,
		},
		GeneratePublishTopic: func(param cqrs.GenerateEventPublishTopicParams) (string, error) {
			return param.EventName, nil
		},
	})
}
