package service

import (
	"github.com/ch0c0-msk/wb-tech-L0/pkg/model"
	"github.com/ch0c0-msk/wb-tech-L0/pkg/repository"
)

type Order interface {
	GetOrder(id string) (model.Order, error)
	AddOrder(order model.Order) error
	AddErrorOrder(order model.Order) error
}

type Service struct {
	Order
}

func NewService(cacheRepo *repository.CacheRepository, dbRepo *repository.DbRepository) *Service {
	return &Service{
		Order: NewOrderService(*cacheRepo, *dbRepo),
	}
}
