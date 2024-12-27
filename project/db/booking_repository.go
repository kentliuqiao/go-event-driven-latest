package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) BookingRepository {
	return BookingRepository{db: db}
}

func (r BookingRepository) Add(ctx context.Context, booking entities.Booking) error {
	_, err := r.db.Exec(
		`INSERT INTO bookings (booking_id, show_id, customer_email, number_of_tickets)
VALUES ($1, $2, $3, $4)
ON CONFLICT (booking_id) DO NOTHING
`, booking.BookingID, booking.ShowID, booking.CustomerEmail, booking.NumberOfTickets)
	if err != nil {
		return fmt.Errorf("could not add booking: %w", err)
	}

	return nil
}
