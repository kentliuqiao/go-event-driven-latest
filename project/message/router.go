package message

import (
	"tickets/message/event"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

const brokenMessageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

func NewWatermillEventProcessor(handler event.Handler, rdb *redis.Client, logger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	useMiddlewares(router, logger)

	marshaler := cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}
	ep, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(param cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return param.EventName, nil
			},
			SubscriberConstructor: func(param cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return redisstream.NewSubscriber(redisstream.SubscriberConfig{
					Client:        rdb,
					ConsumerGroup: "svc-tockets." + param.HandlerName,
				}, logger)
			},
			Marshaler: marshaler,
			Logger:    logger,
		},
	)
	if err != nil {
		panic(err)
	}

	ep.AddHandlers(cqrs.NewEventHandler(
		"issue_receipt",
		handler.IssueReceipt,
	))
	ep.AddHandlers(cqrs.NewEventHandler(
		"append_to_tracker",
		handler.AppendToTracker,
	))
	ep.AddHandlers(cqrs.NewEventHandler(
		"cancel_ticket",
		handler.CancelTicket,
	))
	ep.AddHandlers(cqrs.NewEventHandler(
		"store_ticket",
		handler.StoreTicket,
	))
	ep.AddHandlers(cqrs.NewEventHandler(
		"remove_ticket",
		handler.RemoveCanceledTicket,
	))
	ep.AddHandlers(cqrs.NewEventHandler(
		"generate_ticket",
		handler.PrintTicket,
	))

	return router
}
