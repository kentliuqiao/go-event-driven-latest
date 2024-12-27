package main

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func RegisterEventHandlers(
	sub message.Subscriber,
	router *message.Router,
	handlers []cqrs.EventHandler,
	logger watermill.LoggerAdapter,
) error {
	ep, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(param cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return param.EventName, nil
			},
			SubscriberConstructor: func(param cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return sub, nil
			},
			Marshaler: cqrs.JSONMarshaler{
				GenerateName: cqrs.StructName,
			},
			Logger: logger,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create event processor: %w", err)
	}
	err = ep.AddHandlers(handlers...)
	if err != nil {
		return fmt.Errorf("could not add handlers to event processor: %w", err)
	}

	return nil
}
