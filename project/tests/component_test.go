package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"tickets/api"
	"tickets/entities"
	"tickets/message"
	"tickets/service"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponent(t *testing.T) {
	postgresClient := message.NewPostgresqlClient(os.Getenv("POSTGRES_URL"))
	defer postgresClient.Close()

	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spreadsheetsService := &api.SpreadsheetsMock{}
	receiptsService := &api.ReceiptsMock{}
	fileMock := &api.FileMock{}

	go func() {
		svc := service.New(
			postgresClient,
			redisClient,
			spreadsheetsService,
			receiptsService,
			fileMock,
		)
		assert.NoError(t, svc.Run(ctx))
	}()

	waitForHttpServer(t)

	ticket := TicketStatus{
		TicketID: uuid.NewString(),
		Status:   "confirmed",
		Price: Money{
			Amount:   "100",
			Currency: "USD",
		},
		Email:     "a@a.com",
		BookingID: "asdf",
	}

	sendTicketsStatus(t, TicketsStatusRequest{
		Tickets: []TicketStatus{ticket},
	})

	assertReceiptForTicketIssued(t, receiptsService, ticket)
	assertTicketPrinted(t, fileMock, ticket)
	assertTicketStored(t, postgresClient, ticket)

	sendTicketsStatus(t, TicketsStatusRequest{
		Tickets: []TicketStatus{ticket},
	})

	assertSheetRowAppended(t, spreadsheetsService, TicketStatus{
		TicketID: uuid.NewString(),
		Status:   "canceled",
		Price: Money{
			Amount:   "100",
			Currency: "USD",
		},
		Email:     "a@a.com",
		BookingID: "asdf",
	})
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}

func assertSheetRowAppended(t *testing.T, spreadsheetsService *api.SpreadsheetsMock, event TicketStatus) {
	assert.EventuallyWithT(
		t,
		func(collectT *assert.CollectT) {
			appendedRows := len(spreadsheetsService.Rows)
			t.Log("appended rows", appendedRows)

			assert.Greater(collectT, appendedRows, 0, "no rows appended")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	var row []string
	var ok bool
LOOP:
	for sheetName, appendedRow := range spreadsheetsService.Rows {
		if sheetName != "tickets-to-refund" {
			continue
		}
		for _, appendedRow := range appendedRow {
			if appendedRow[0] != event.TicketID {
				continue
			}
			row = appendedRow
			ok = true
			break LOOP
		}
	}
	require.True(t, ok, "row not appended")

	assert.Equal(t, event.TicketID, row[0])
	assert.Equal(t, event.Price.Amount, row[2])
	assert.Equal(t, event.Price.Currency, row[3])
}

func assertReceiptForTicketIssued(t *testing.T, receiptsService *api.ReceiptsMock, ticket TicketStatus) {
	assert.EventuallyWithT(
		t,
		func(collectT *assert.CollectT) {
			issuedReceipts := len(receiptsService.IssuedReceipts)
			t.Log("issued receipts", issuedReceipts)

			assert.Greater(collectT, issuedReceipts, 0, "no receipts issued")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	var receipt entities.IssueReceiptRequest
	var ok bool
	for _, issuedReceipt := range receiptsService.IssuedReceipts {
		if issuedReceipt.TicketID != ticket.TicketID {
			continue
		}
		receipt = issuedReceipt
		ok = true
		break
	}
	require.Truef(t, ok, "receipt for ticket %s not found", ticket.TicketID)

	assert.Equal(t, ticket.TicketID, receipt.TicketID)
	assert.Equal(t, ticket.Price.Amount, receipt.Price.Amount)
	assert.Equal(t, ticket.Price.Currency, receipt.Price.Currency)
}

func assertTicketPrinted(t *testing.T, fileMock *api.FileMock, ticket TicketStatus) {
	assert.EventuallyWithT(
		t,
		func(collectT *assert.CollectT) {
			printedTickets := len(fileMock.Files)
			t.Log("printed tickets", printedTickets)

			assert.Greater(collectT, printedTickets, 0, "no tickets printed")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	fileMock.DownloadFile(context.Background(), entities.DownloadFileRequest{
		FileID: ticket.TicketID,
	})

	var ticketID string
	var ok bool
	for _, printedTicket := range fileMock.Files {
		if printedTicket != ticket.TicketID {
			continue
		}
		ticketID = printedTicket
		ok = true
		break
	}
	require.Truef(t, ok, "ticket %s not printed", ticket.TicketID)

	assert.Equal(t, ticket.TicketID, ticketID)
}

func assertTicketStored(t *testing.T, postgresClient *sqlx.DB, ticket TicketStatus) {
	rows, err := postgresClient.Query(`SELECT * FROM tickets WHERE ticket_id = $1`, ticket.TicketID)
	require.NoError(t, err)
	defer rows.Close()

	var count int
	for rows.Next() {
		count++
		var (
			ticketID  string
			status    string
			price     string
			email     string
			bookingID string
		)

		err := rows.Scan(&ticketID, &status, &price, &email, &bookingID)
		require.NoError(t, err)

		assert.Equal(t, ticket.TicketID, ticketID)
		assert.Equal(t, ticket.Status, status)
		assert.Equal(t, ticket.Price.Amount, price)
		assert.Equal(t, ticket.Email, email)
		assert.Equal(t, ticket.BookingID, bookingID)
	}

	assert.Equal(t, 1, count)
}

type TicketsStatusRequest struct {
	Tickets []TicketStatus `json:"tickets"`
}

type TicketStatus struct {
	TicketID  string `json:"ticket_id"`
	Status    string `json:"status"`
	Price     Money  `json:"price"`
	Email     string `json:"email"`
	BookingID string `json:"booking_id"`
}

type Money struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func sendTicketsStatus(t *testing.T, req TicketsStatusRequest) {
	t.Helper()

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	correlationID := shortuuid.New()

	ticketIDs := make([]string, 0, len(req.Tickets))
	for _, ticket := range req.Tickets {
		ticketIDs = append(ticketIDs, ticket.TicketID)
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/tickets-status",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)

	httpReq.Header.Set("Correlation-ID", correlationID)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotency-Key", shortuuid.New())

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
