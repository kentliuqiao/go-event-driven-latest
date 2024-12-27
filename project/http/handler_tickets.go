package http

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ticketsStatusRequest struct {
	Tickets []ticketStatusRequest `json:"tickets"`
}

type ticketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
	BookingID     string         `json:"booking_id"`
}

func (h Handler) PostTicketsStatus(c echo.Context) error {
	idempotencyKey := c.Request().Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		return fmt.Errorf("missing Idempotency-Key header")
	}

	var request ticketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		if ticket.Status == "confirmed" {
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewEventHeaderWithIdempotencyKey(idempotencyKey),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			err = h.eb.Publish(c.Request().Context(), event)
			if err != nil {
				return err
			}
		} else if ticket.Status == "canceled" {
			event := entities.TicketBookingCanceled{
				Header:        entities.NewEventHeaderWithIdempotencyKey(idempotencyKey),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			err = h.eb.Publish(c.Request().Context(), event)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h Handler) GetTickets(c echo.Context) error {
	tickets, err := h.ticketRepo.FindAll(c.Request().Context())
	if err != nil {
		return fmt.Errorf("failed to find all tickets: %w", err)
	}

	return c.JSON(http.StatusOK, tickets)
}

func (h Handler) PostBookTickets(c echo.Context) error {
	var req entities.Booking
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	req.BookingID = uuid.NewString()
	err = h.bookingRepo.Add(c.Request().Context(), req)
	if err != nil {
		return fmt.Errorf("failed to add booking: %w", err)
	}

	return c.JSON(http.StatusCreated, req)
}
