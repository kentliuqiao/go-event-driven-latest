package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SubscribeForMessages(db *sqlx.DB, topic string, logger watermill.LoggerAdapter) (<-chan *message.Message, error) {
	subsriber, err := watermillSQL.NewSubscriber(db, watermillSQL.SubscriberConfig{
		SchemaAdapter:  watermillSQL.DefaultPostgreSQLSchema{},
		OffsetsAdapter: watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create sql subscriber: %w", err)
	}
	err = subsriber.SubscribeInitialize(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize subscription: %w", err)
	}

	return subsriber.Subscribe(context.Background(), topic)
}
