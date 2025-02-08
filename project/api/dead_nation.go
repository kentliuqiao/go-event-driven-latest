package api

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/dead_nation"
)

type DeadNationClient struct {
	clients *clients.Clients
}

func NewDeadNationAPI(clients *clients.Clients) *DeadNationClient {
	return &DeadNationClient{clients: clients}
}

func (d DeadNationClient) PostTicketBookingWithResponse(ctx context.Context, booking entities.DeadNationBooking) (*dead_nation.PostTicketBookingResponse, error) {
	resp, err := d.clients.DeadNation.PostTicketBookingWithResponse(ctx, dead_nation.PostTicketBookingRequest{
		BookingId:       booking.BookingID,
		CustomerAddress: booking.CustomerEmail,
		EventId:         booking.DeadNationEventID,
		NumberOfTickets: booking.NumberOfTickets,
	})
	if err != nil {
		return nil, fmt.Errorf("could not post ticket booking to dead nation: %w", err)
	}

	return resp, nil
}
