package http

import (
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h Handler) PostShows(c echo.Context) error {
	var shows entities.Show
	err := c.Bind(&shows)
	if err != nil {
		return err
	}

	show := entities.Show{
		ShowID:          uuid.NewString(),
		DeadNationID:    shows.DeadNationID,
		NumberOfTickets: shows.NumberOfTickets,
		StartTime:       shows.StartTime,
		Title:           shows.Title,
		Venue:           shows.Venue,
	}

	err = h.showRepo.AddShow(c.Request().Context(), show)
	if err != nil {
		return fmt.Errorf("failed to add show: %w", err)
	}

	return c.JSON(http.StatusCreated, show)
}
