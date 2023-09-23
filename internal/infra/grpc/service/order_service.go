package service

import (
	"context"

	"github.com/gbgomes/GoExpert/CleanArch/internal/infra/grpc/pb"
	"github.com/gbgomes/GoExpert/CleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase   usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(context.Context, *pb.Blank) (*pb.OrderList, error) {
	dto := usecase.OrderListInputDTO{}
	output, err := s.ListOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}

	outs := []*pb.ListOrderResponse{}
	for _, o := range output {
		out := pb.ListOrderResponse{
			Id:         o.ID,
			Price:      float32(o.Price),
			Tax:        float32(o.Tax),
			FinalPrice: float32(o.FinalPrice),
		}
		outs = append(outs, &out)
	}

	return &pb.OrderList{Orders: outs}, nil
}
