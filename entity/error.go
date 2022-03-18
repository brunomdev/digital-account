package entity

import "errors"

var ErrNotFound = errors.New("not found")
var ErrInvalidAmount = errors.New("invalid amount")
var ErrInsufficientCreditLimit = errors.New("available credit limit is insufficient")
