package models

import "errors"

type Order struct {
	ID        int
	UserID    int
	OrderID   string
	Accrual   float32
	Status    string
	CreatedAt string
}

var ErrOrderExists = errors.New("order exists")
var ErrOrderNotFound = errors.New("order not found")
var ErrOrderInvalid = errors.New("order invalid")
var ErrOrderUploadedSameUser = errors.New("order already loaded same user")
var ErrOrderUploadedAnotherUser = errors.New("order loaded another user")
