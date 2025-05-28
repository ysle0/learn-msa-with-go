package main

import (
	"github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
)

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, i := range items {
		if i.ID == "" {
			return common.ErrItemIdRequired
		}

		if i.Quantity <= 0 {
			return common.ErrItemQuantityRequired
		}
	}

	return nil
}
