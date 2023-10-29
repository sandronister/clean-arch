package entity

type OrderRespositoryInterface interface {
	Save(order *Order) error
}
