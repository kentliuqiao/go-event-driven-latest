package event

import (
	"context"
	"tickets/db"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	ticketRepo          db.TicketRepository
	spreadsheetsService SpreadsheetsAPI
	receiptsService     ReceiptsService
	filesAPI            FilesAPI
	eventBus            *cqrs.EventBus
}

func NewHandler(
	dbC *sqlx.DB,
	spreadsheetsService SpreadsheetsAPI,
	receiptsService ReceiptsService,
	filesAPI FilesAPI,
	eb *cqrs.EventBus,
) Handler {
	if spreadsheetsService == nil {
		panic("missing spreadsheetsService")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		ticketRepo:          db.NewTicketRepository(dbC),
		spreadsheetsService: spreadsheetsService,
		receiptsService:     receiptsService,
		filesAPI:            filesAPI,
		eventBus:            eb,
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type SpreadsheetsAPIMock struct {
}

func (m SpreadsheetsAPIMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	return nil
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type ReceiptsServiceMock struct {
}

func (m *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	return entities.IssueReceiptResponse{}, nil
}

type FilesAPI interface {
	UploadFile(ctx context.Context, req entities.GenerateFileRequest) error
	DownloadFile(ctx context.Context, req entities.DownloadFileRequest) ([]byte, error)
}
