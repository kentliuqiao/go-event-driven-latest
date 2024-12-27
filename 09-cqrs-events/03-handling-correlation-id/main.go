package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/lithammer/shortuuid/v3"
)

type CorrelationPublisherDecorator struct {
	message.Publisher
}

func (c CorrelationPublisherDecorator) Publish(topic string, messages ...*message.Message) error {
	for i := range messages {
		// if correlation_id is already set, let's not override
		if messages[i].Metadata.Get("correlation_id") != "" {
			continue
		}

		// correlation_id as const
		messages[i].Metadata.Set("correlation_id", CorrelationIDFromContext(messages[i].Context()))
	}

	return c.Publisher.Publish(topic, messages...)
}

type ctxKey int

const (
	correlationIDKey ctxKey = iota
)

func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	v, ok := ctx.Value(correlationIDKey).(string)
	if ok {
		return v
	}

	// add "gen_" prefix to distinguish generated correlation IDs from correlation IDs passed by the client
	// it's useful to detect if correlation ID was not passed properly
	return "gen_" + shortuuid.New()
}
