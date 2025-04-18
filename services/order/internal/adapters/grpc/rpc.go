package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hollowdll/go-grpc-microservices/services/order/internal/application/core/domain"
	"github.com/hollowdll/grpc-microservices-proto/gen/golang/orderpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	createOrderName    string = "CreateOrder"
	createOrderTimeout        = time.Second * 5
)

func (a *Adapter) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (res *orderpb.CreateOrderResponse, err error) {
	log.Printf("call RPC %s: request = %v", createOrderName, req)
	defer func() {
		if err != nil {
			log.Printf("RPC %s failed: request = %v; error = %v", createOrderName, req, err)
		} else {
			log.Printf("RPC %s success: request = %v; response = %v", createOrderName, req, res)
		}
	}()

	// use timeout so the order operation won't wait forever
	ctx, cancel := context.WithTimeout(ctx, createOrderTimeout)
	defer cancel()

	var orderItems []*domain.OrderItem
	for _, orderItem := range req.OrderItems {
		orderItems = append(orderItems, &domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			Quantity:    orderItem.Quantity,
		})
	}
	newOrder, err := domain.NewOrder(req.CustomerId, orderItems)
	if err != nil {
		log.Printf("RPC %s: failed to create new order: %v", createOrderName, err)
		return nil, status.Error(codes.Internal, "order creation failed")
	}

	result, err := a.api.CreateOrder(ctx, newOrder)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("order creation failed: %v", err))
	}

	return &orderpb.CreateOrderResponse{OrderId: result.ID}, nil
}
