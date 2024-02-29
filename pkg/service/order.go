package service

import (
	"github.com/ch0c0-msk/wb-tech-L0/pkg/model"
	"github.com/ch0c0-msk/wb-tech-L0/pkg/repository"
)

type OrderService struct {
	cacheRepo repository.CacheRepository
	dbRepo    repository.DbRepository
}

func NewOrderService(cacheRepo repository.CacheRepository, dbRepo repository.DbRepository) *OrderService {
	return &OrderService{
		cacheRepo: cacheRepo,
		dbRepo:    dbRepo,
	}
}

func (o *OrderService) GetOrder(id string) (model.Order, error) {
	return o.cacheRepo.GetOrder(id)
}

func (o *OrderService) AddOrder(order model.Order) error {
	if err := o.dbRepo.AddOrder(order); err != nil {
		return err
	}
	return o.cacheRepo.AddOrder(order)
}
