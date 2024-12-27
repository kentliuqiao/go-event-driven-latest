package event

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) RemoveCanceledTicket(ctx context.Context, event *entities.TicketBookingCanceled) error {
	log.FromContext(ctx).Info("Removing ticket")

	return h.ticketRepo.Remove(ctx, event.TicketID)
}
