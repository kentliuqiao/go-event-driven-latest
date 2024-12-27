package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type ShowRepository struct {
	db *sqlx.DB
}

func NewShowRepository(db *sqlx.DB) ShowRepository {
	if db == nil {
		panic("missing db")
	}
	return ShowRepository{db: db}
}

func (r ShowRepository) Add(ctx context.Context, show entities.Show) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO shows (show_id, dead_nation_id, number_of_tickets, start_time, title, venue)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (show_id) DO NOTHING
`, show.ShowID, show.DeadNationID, show.NumberOfTickets, show.StartTime, show.Title, show.Venue)
	if err != nil {
		return fmt.Errorf("could not add show: %w", err)
	}

	return nil
}

func (r ShowRepository) Remove(showID string) error {
	_, err := r.db.Exec(
		`DELETE FROM shows WHERE show_id = $1`,
		showID)
	if err != nil {
		return fmt.Errorf("could not delete show: %w", err)
	}

	return nil
}

func (r ShowRepository) FindAll() ([]entities.Show, error) {
	var shows []entities.Show

	err := r.db.Select(
		&shows,
		`SELECT show_id, dead_nation_id, number_of_tickets, start_time, title, venue FROM shows`,
	)
	if err != nil {
		return nil, fmt.Errorf("could not list shows: %w", err)
	}

	return shows, nil
}
