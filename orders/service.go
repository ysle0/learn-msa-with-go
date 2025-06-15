package main

import (
	"context"
	"log"

	"github.com/ysle0/omsv2/common"
	pb "github.com/ysle0/omsv2/common/api"
)

type service struct {
	store OrdersStore
}

func NewService(store OrdersStore) *service {
	return &service{store}
}

func (s *service) CreateOrder(ctx context.Context, r *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrder(ctx, r)
	if err != nil {
		return nil, err
	}

	log.Printf("order validated: %v\n", items)
	o := &pb.Order{
		ID:         "42",
		CustomerID: r.CustomerID,
		Status:     "pending",
		Items:      items,
	}

	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, r *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(r.Items) == 0 {
		log.Printf("no items in request")
		return nil, common.ErrNoItems
	}

	log.Printf("starting merge with items: %v\n", r.Items)
	mergedItems := mergeItemsQuantities(r.Items)
	log.Printf("merged items: %v\n", mergedItems)

	// validate with the stock service
	var itemsWithPrice []*pb.Item
	for _, item := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			ID:       item.ID,
			PriceID:  "price_1RaIJzFM42raIhkYqptQbQNR",
			Quantity: item.Quantity,
		})
	}

	return itemsWithPrice, nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	var merged []*pb.ItemsWithQuantity
	for _, item := range items {
		found := false

		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
