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
	Withdraw  float32
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

type GetAllWithdraw struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

var ErrOrderExists = errors.New("order exists")
var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrOrderInvalid = errors.New("order invalid")
var ErrOrderUploadedSameUser = errors.New("order already loaded same user")
var ErrOrderUploadedAnotherUser = errors.New("order loaded another user")
