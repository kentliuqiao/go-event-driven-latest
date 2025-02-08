package api

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients/dead_nation"
)

type DeadNationClientMock struct{}

func (m DeadNationClientMock) PostTicketBookingWithResponse(ctx context.Context, booking entities.DeadNationBooking) (*dead_nation.PostTicketBookingResponse, error) {
	return nil, nil
}
