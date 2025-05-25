package common

import "errors"

var (
	ErrNoItems = errors.New("items is empty")
	ErrItemIdRequired = errors.New("item id is required")
	ErrItemQuantityRequired = errors.New("item quantity is required")
)
