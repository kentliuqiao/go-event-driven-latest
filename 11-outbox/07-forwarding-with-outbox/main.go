package main

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	watermillSQL "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	_ "github.com/lib/pq"
)

func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {
	sub, err := watermillSQL.NewSubscriber(db, watermillSQL.SubscriberConfig{
		SchemaAdapter:  watermillSQL.DefaultPostgreSQLSchema{},
		OffsetsAdapter: watermillSQL.DefaultPostgreSQLOffsetsAdapter{},
	}, logger)
	if err != nil {
		return fmt.Errorf("failed to create sql subscriber: %w", err)
	}
	err = sub.SubscribeInitialize(outboxTopic)
	if err != nil {
		return fmt.Errorf("failed to initialize subscription: %w", err)
	}

	redisPub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		return fmt.Errorf("failed to create redis publisher: %w", err)
	}

	fw, err := forwarder.NewForwarder(sub, redisPub, logger, forwarder.Config{
		ForwarderTopic: outboxTopic,
	})
	if err != nil {
		return fmt.Errorf("failed to create forwarder: %w", err)
	}

	go func() {
		err := fw.Run(context.Background())
		if err != nil {
			panic(fmt.Errorf("failed to run forwarder: %w", err))
		}
	}()

	<-fw.Running()

	return nil
}
