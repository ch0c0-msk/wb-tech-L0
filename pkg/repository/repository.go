package repository

import (
	"database/sql"

	"github.com/ch0c0-msk/wb-tech-L0/pkg/model"
)

type CacheStore interface {
	GetOrder(id string) (model.Order, error)
	AddOrder(order model.Order) error
}

type DbStore interface {
	AddOrder(order model.Order) error
	AddErrorOrder(order model.Order) error
	RestoreCache() (map[string]model.Order, error)
}

type CacheRepository struct {
	CacheStore
}

type DbRepository struct {
	DbStore
}

func NewCacheRepository(db *sql.DB) (*CacheRepository, error) {
	cacheDB, err := NewCacheDB(*NewOrderSql(db))
	if err != nil {
		return nil, err
	}
	return &CacheRepository{CacheStore: cacheDB}, nil
}

func NewDbRepository(db *sql.DB) *DbRepository {
	return &DbRepository{
		DbStore: NewOrderSql(db),
	}
}
