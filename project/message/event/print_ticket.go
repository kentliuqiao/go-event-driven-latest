package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) PrintTicket(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("Printing ticket")

	ticketHTML := `
		<html>
			<head>
				<title>Ticket</title>
			</head>
			<body>
				<h1>Ticket ` + event.TicketID + `</h1>
				<p>Price: ` + event.Price.Amount + ` ` + event.Price.Currency + `</p>	
			</body>
		</html>
`

	ticketFile := event.TicketID + "-ticket.html"

	err := h.filesAPI.UploadFile(ctx, entities.GenerateFileRequest{
		FileID:      ticketFile,
		FileContent: ticketHTML,
	})
	if err != nil {
		return fmt.Errorf("failed to upload ticket file: %w", err)
	}

	err = h.eventBus.Publish(ctx, entities.TicketPrinted{
		Header:   entities.NewEventHeader(),
		TicketID: event.TicketID,
		FileName: ticketFile,
	})
	if err != nil {
		return fmt.Errorf("failed to publish TicketPrinted event: %w", err)
	}

	return nil
}
