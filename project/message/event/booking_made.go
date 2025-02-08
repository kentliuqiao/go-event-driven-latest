package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) BookingMade(ctx context.Context, event *entities.BookingMade) error {
	log.FromContext(ctx).Info("Booking ticket in Dead Nation")

	show, err := h.showsRepository.ShowByID(ctx, event.ShowId)
	if err != nil {
		return fmt.Errorf("failed to get show: %w", err)
	}

	resp, err := h.deadNationClient.PostTicketBookingWithResponse(ctx, entities.DeadNationBooking{
		BookingID:         event.BookingID,
		NumberOfTickets:   event.NumberOfTickets,
		CustomerEmail:     event.CustomerEmail,
		DeadNationEventID: show.DeadNationID,
	})
	if err != nil {
		return fmt.Errorf("could not post ticket booking: %w", err)
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return nil
}
