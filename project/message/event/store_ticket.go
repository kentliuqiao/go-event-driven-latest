package event

import (
	"context"
	"tickets/entities"
)

func (h Handler) StoreTicket(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	err := h.ticketRepo.Add(ctx, entities.Ticket{
		TicketID:      event.TicketID,
		Price:         event.Price,
		CustomerEmail: event.CustomerEmail,
	})
	if err != nil {
		return err
	}

	return nil
}
