package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"tickets/api"
	"tickets/message"
	"tickets/service"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	apiClients, err := clients.NewClients(
		os.Getenv("GATEWAY_ADDR"),
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Correlation-ID", log.CorrelationIDFromContext(ctx))
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	postgresClient := message.NewPostgresqlClient(os.Getenv("POSTGRES_URL"))
	defer postgresClient.Close()

	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	spreadsheetsService := api.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := api.NewReceiptsServiceClient(apiClients)
	filesAPI := api.NewFilesAPIClient(apiClients)
	deadNationAPI := api.NewDeadNationAPI(apiClients)

	err = service.New(
		postgresClient,
		redisClient,
		spreadsheetsService,
		receiptsService,
		filesAPI,
		deadNationAPI,
	).Run(ctx)
	if err != nil {
		panic(err)
	}
}
