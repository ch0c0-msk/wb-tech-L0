package repository

import (
	"errors"
	"fmt"
	"sync"

	"github.com/ch0c0-msk/wb-tech-L0/pkg/model"
)

type cacheDB struct {
	sync.Mutex
	orders map[string]model.Order
}

func NewCacheDB(orderSql OrderSql) (*cacheDB, error) {
	orders, err := orderSql.RestoreCache()
	if err != nil {
		return nil, fmt.Errorf("failed to restore cache: %s", err.Error())
	}
	return &cacheDB{orders: orders}, nil
}

func (db *cacheDB) GetOrder(id string) (model.Order, error) {
	db.Mutex.Lock()
	order, ok := db.orders[id]
	db.Mutex.Unlock()
	if !ok {
		return order, errors.New("order didnt find in cache")
	}
	return order, nil
}

func (db *cacheDB) AddOrder(order model.Order) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	_, exist := db.orders[order.Id]
	if exist {
		return errors.New("order already exist")
	}
	db.orders[order.Id] = order
	return nil
}
