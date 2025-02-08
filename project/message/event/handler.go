package event

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients/dead_nation"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	spreadsheetsService SpreadsheetsAPI
	receiptsService     ReceiptsService
	filesAPI            FilesAPI
	eventBus            *cqrs.EventBus
	ticketRepo          TicketsRepository
	showsRepository     ShowsRepository
	deadNationClient    DeadNationClient
}

func NewHandler(
	dbC *sqlx.DB,
	spreadsheetsService SpreadsheetsAPI,
	receiptsService ReceiptsService,
	filesAPI FilesAPI,
	eb *cqrs.EventBus,
	ticketRepo TicketsRepository,
	showsRepo ShowsRepository,
	deadNationClient DeadNationClient,
) Handler {
	if spreadsheetsService == nil {
		panic("missing spreadsheetsService")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}

	return Handler{
		spreadsheetsService: spreadsheetsService,
		receiptsService:     receiptsService,
		filesAPI:            filesAPI,
		eventBus:            eb,
		ticketRepo:          ticketRepo,
		showsRepository:     showsRepo,
		deadNationClient:    deadNationClient,
	}
}

type SpreadsheetsAPIMock struct {
}

func (m SpreadsheetsAPIMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	return nil
}

type ReceiptsServiceMock struct {
}

func (m *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	return entities.IssueReceiptResponse{}, nil
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type FilesAPI interface {
	UploadFile(ctx context.Context, req entities.GenerateFileRequest) error
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket entities.Ticket) error
	Remove(ctx context.Context, ticketID string) error
}

type ShowsRepository interface {
	ShowByID(ctx context.Context, showID uuid.UUID) (entities.Show, error)
}

type DeadNationClient interface {
	PostTicketBookingWithResponse(ctx context.Context, booking entities.DeadNationBooking) (*dead_nation.PostTicketBookingResponse, error)
}
