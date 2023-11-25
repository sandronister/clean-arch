package usecase

import "github.com/sandronister/clean-arch/internal/entity"

type ListOrderUseCase struct {
	OrderRepository entity.OrderRespositoryInterface
}

func NewListOrderUseCase(OrderRepository entity.OrderRespositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := l.OrderRepository.List()

	if err != nil {
		return []OrderOutputDTO{}, err
	}

	outputDTO := []OrderOutputDTO{}

	for _, order := range orders {
		outputDTO = append(outputDTO, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return outputDTO, nil
}
