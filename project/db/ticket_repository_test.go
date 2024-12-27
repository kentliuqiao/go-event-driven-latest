package db_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"tickets/db"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestTicketRepository_Add(t *testing.T) {
	type args struct {
		ctx    context.Context
		ticket entities.Ticket
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should add ticket",
			args: args{
				ctx: context.Background(),
				ticket: entities.Ticket{
					TicketID:      uuid.New().String(),
					Price:         entities.Money{Amount: "100", Currency: "USD"},
					CustomerEmail: "x@x.com",
				},
			},
			wantErr: false,
		},
	}
	err := db.InitializeDbSchema(getDb())
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := db.NewTicketRepository(getDb())
			if err := r.Add(tt.args.ctx, tt.args.ticket); (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			tickets, err := r.FindAll(tt.args.ctx)
			if err != nil {
				t.Errorf("TicketRepository.FindAll() error = %v", err)
			}
			for _, ticket := range tickets {
				if ticket.TicketID == tt.args.ticket.TicketID {
					return
				}
			}
			t.Errorf("TicketRepository.FindAll() = %v, want %v", tickets, tt.args.ticket)
		})
	}
}

func TestTicketRepository_Add_Idempotent(t *testing.T) {
	type args struct {
		ctx    context.Context
		ticket entities.Ticket
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should add ticket only once",
			args: args{
				ctx: context.Background(),
				ticket: entities.Ticket{
					TicketID:      uuid.New().String(),
					Price:         entities.Money{Amount: "100", Currency: "USD"},
					CustomerEmail: "x@x.com",
				},
			},
			wantErr: false,
		},
	}
	err := db.InitializeDbSchema(getDb())
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := db.NewTicketRepository(getDb())
			if err := r.Add(tt.args.ctx, tt.args.ticket); (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := r.Add(tt.args.ctx, tt.args.ticket); (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			tickets, err := r.FindAll(tt.args.ctx)
			if err != nil {
				t.Errorf("TicketRepository.FindAll() error = %v", err)
			}
			var added bool
			for _, ticket := range tickets {
				if ticket.TicketID == tt.args.ticket.TicketID {
					added = true
				}
			}
			if !added {
				t.Errorf("TicketRepository.FindAll() = %v, want %v", tickets, tt.args.ticket)
			}
			if err := r.Add(tt.args.ctx, tt.args.ticket); (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			tickets, err = r.FindAll(tt.args.ctx)
			if err != nil {
				t.Errorf("TicketRepository.FindAll() error = %v", err)
			}
			var count int
			for _, ticket := range tickets {
				if ticket.TicketID == tt.args.ticket.TicketID {
					count++
				}
			}
			if count != 1 {
				t.Errorf("TicketRepository.FindAll() = %v, want %v", tickets, tt.args.ticket)
			}
			if err := r.Remove(tt.args.ctx, tt.args.ticket.TicketID); (err != nil) != tt.wantErr {
				t.Errorf("TicketRepository.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
			tickets, err = r.FindAll(tt.args.ctx)
			if err != nil {
				t.Errorf("TicketRepository.FindAll() error = %v", err)
			}
			for _, ticket := range tickets {
				if ticket.TicketID == tt.args.ticket.TicketID {
					t.Errorf("TicketRepository.FindAll() = %v, want %v", tickets, tt.args.ticket)
				}
			}
			t.Log("tickets", tickets)
		})
	}
}

var dbc *sqlx.DB
var getDbOnce sync.Once

func getDb() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		dbc, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}
	})
	return dbc
}
