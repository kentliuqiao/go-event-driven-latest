package http

import (
	"net/http"
	"tickets/db"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(eb *cqrs.EventBus, spreadsheetsAPIClient SpreadsheetsAPI, dbc *sqlx.DB) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eb:                    eb,
		spreadsheetsAPIClient: spreadsheetsAPIClient,
		ticketRepo:            db.NewTicketRepository(dbc),
		showRepo:              db.NewShowRepository(dbc),
		bookingRepo:           db.NewBookingRepository(dbc),
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)

	e.GET("/tickets", handler.GetTickets)

	e.POST("/shows", handler.PostShows)

	e.POST("/book-tickets", handler.PostBookTickets)

	return e
}
