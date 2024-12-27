package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type TicketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	if db == nil {
		panic("missing db")
	}
	return TicketRepository{db: db}
}

func (r TicketRepository) Add(ctx context.Context, ticket entities.Ticket) error {
	_, err := r.db.ExecContext(ctx, `
INSERT INTO tickets (ticket_id, price_amount, price_currency, customer_email)
VALUES ($1, $2, $3, $4)
ON CONFLICT (ticket_id) DO NOTHING
`, ticket.TicketID, ticket.Price.Amount, ticket.Price.Currency, ticket.CustomerEmail)
	if err != nil {
		return fmt.Errorf("could not add ticket: %w", err)
	}

	return nil
}

func (r TicketRepository) Remove(ctx context.Context, ticketID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM tickets WHERE ticket_id = $1`,
		ticketID)
	if err != nil {
		return fmt.Errorf("could not delete ticket: %w", err)
	}

	return nil
}

func (r TicketRepository) FindAll(ctx context.Context) ([]entities.Ticket, error) {
	var tickets []entities.Ticket

	err := r.db.SelectContext(
		ctx,
		&tickets,
		`SELECT ticket_id, price_amount as "price.amount", price_currency as "price.currency", customer_email FROM tickets`,
	)
	if err != nil {
		return nil, fmt.Errorf("could not list tickets: %w", err)
	}

	return tickets, nil
}
