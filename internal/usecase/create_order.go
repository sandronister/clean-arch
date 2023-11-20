package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sandronister/clean-arch/internal/entity"
	"github.com/sandronister/clean-arch/pkg/events"
)

type OrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRespositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository entity.OrderRespositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input OrderInputDTO) (OrderOutputDTO, error) {

	order, err := entity.NewOrder(uuid.NewString(), input.Price, input.Tax)

	if err != nil {
		return OrderOutputDTO{}, err
	}

	err = order.CalculateFinalPrice()

	if err != nil {
		return OrderOutputDTO{}, nil
	}

	if err = c.OrderRepository.Save(order); err != nil {
		fmt.Println(err)
		return OrderOutputDTO{}, nil
	}

	dto := OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil

}
