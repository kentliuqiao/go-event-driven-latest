package service

import (
	"context"
	stdHTTP "net/http"
	"tickets/db"
	ticketsHttp "tickets/http"
	"tickets/message"
	"tickets/message/event"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.Init(logrus.InfoLevel)
}

type Service struct {
	db              *sqlx.DB
	watermillRouter *watermillMessage.Router
	echoRouter      *echo.Echo
}

func New(
	db *sqlx.DB,
	redisClient *redis.Client,
	spreadsheetsService event.SpreadsheetsAPI,
	receiptsService event.ReceiptsService,
	filesAPI event.FilesAPI,
) Service {
	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}

	marshaler := cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}
	eventBus, err := cqrs.NewEventBusWithConfig(
		redisPublisher,
		cqrs.EventBusConfig{
			GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
				return params.EventName, nil
			},
			Marshaler: marshaler,
		},
	)
	if err != nil {
		panic(err)
	}

	handler := event.NewHandler(db, spreadsheetsService, receiptsService, filesAPI, eventBus)

	watermillRouter := message.NewWatermillEventProcessor(
		handler,
		redisClient,
		watermillLogger,
	)

	echoRouter := ticketsHttp.NewHttpRouter(
		eventBus,
		spreadsheetsService,
		db,
	)

	return Service{
		db,
		watermillRouter,
		echoRouter,
	}
}

func (s Service) Run(
	ctx context.Context,
) error {
	if err := db.InitializeDbSchema(s.db); err != nil {
		return err
	}

	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")

		if err != nil && err != stdHTTP.ErrServerClosed {
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
