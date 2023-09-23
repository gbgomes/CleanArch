package usecase

import (
	"github.com/gbgomes/GoExpert/CleanArch/internal/entity"
)

type OrderListInputDTO struct {
	Sort   string `json:"sort"`
	Page   int    `json:"page"`
	Offset int    `json:"offset"`
}

type OrderListOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(OrderRepository entity.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute(input OrderListInputDTO) ([]OrderListOutputDTO, error) {
	orders, err := c.OrderRepository.FindAll(input.Page, input.Offset, input.Sort)
	if err != nil {
		return []OrderListOutputDTO{}, err
	}

	out := []OrderListOutputDTO{}
	for _, o := range orders {
		dto := OrderListOutputDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		}
		out = append(out, dto)
	}

	return out, nil
}
