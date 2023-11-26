package services

import (
	"context"

	"github.com/sandronister/clean-arch/internal/infra/grpc/pb"
	"github.com/sandronister/clean-arch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func (s *OrderService) CreaterOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dto := usecase.OrderInputDTO{
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}

	output, err := s.CreateOrderUseCase.Execute(dto)

	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Blank) (*pb.OrderList, error) {
	orders, err := s.ListOrderUseCase.Execute()

	if err != nil {
		return nil, err
	}

	var output = []*pb.OrderResponse{}

	for _, item := range orders {
		output = append(output, &pb.OrderResponse{
			Id:         item.ID,
			Price:      float32(item.Price),
			Tax:        float32(item.Tax),
			FinalPrice: float32(item.FinalPrice),
		})
	}

	return &pb.OrderList{
		Orders: output,
	}, nil
}
