package api

import (
	"context"
	"sync"
	"tickets/entities"
	"time"
)

type ReceiptsMock struct {
	mu             sync.Mutex
	IssuedReceipts []entities.IssueReceiptRequest
}

func (r *ReceiptsMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.IssuedReceipts = append(r.IssuedReceipts, request)

	return entities.IssueReceiptResponse{
		ReceiptNumber: "123456",
		IssuedAt:      time.Now(),
	}, nil
}
