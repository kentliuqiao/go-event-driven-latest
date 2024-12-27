package http

import (
	"context"
	"tickets/db"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	eb                    *cqrs.EventBus
	spreadsheetsAPIClient SpreadsheetsAPI
	ticketRepo            db.TicketRepository
	showRepo              db.ShowRepository
	bookingRepo           db.BookingRepository
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}
