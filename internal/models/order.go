package models

import (
	"errors"
)

const (
	StatusOrderNew        = "NEW"
	StatusOrderProcessing = "PROCESSING"
	StatusOrderInvalid    = "INVALID"
	StatusOrderProcessed  = "PROCESSED"
)

type Order struct {
	ID        int
	UserID    int
	OrderID   string
	Accrual   float32
	Status    string
	CreatedAt string
}

type GetUserOrders struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float32 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
	//Error      string  `json:"error,omitempty"`
}

var ErrOrderExists = errors.New("order exists")
var ErrOrderNotFound = errors.New("order not found")
var ErrOrderInvalid = errors.New("order invalid")
var ErrOrderUploadedSameUser = errors.New("order already loaded same user")
var ErrOrderUploadedAnotherUser = errors.New("order loaded another user")
