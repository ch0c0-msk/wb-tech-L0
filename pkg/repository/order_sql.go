package repository

import (
	"database/sql"
	"fmt"

	"github.com/ch0c0-msk/wb-tech-L0/pkg/model"
)

type OrderSql struct {
	db *sql.DB
}

func NewOrderSql(db *sql.DB) *OrderSql {
	return &OrderSql{db: db}
}

func (o *OrderSql) RestoreCache() (map[string]model.Order, error) {
	orders := make(map[string]model.Order)

	query := fmt.Sprintf(`
		SELECT t1.order_uid, t1.track_number, t1.entry_name, t1.locale, t1.internal_signature, t1.customer_id, t1.delivery_service, t1.shardkey,
			t1.sm_id, t1.date_created, t2.delivery_name, t2.phone, t2.zip, t2.city, t2.delivery_address, t2.region, t2.email, t3.transaction_id,
			t3.request_id, t3.currency, t3.provider_name, t3.amount, t3.payment_dt, t3.bank, t3.delivery_cost, t3.goods_total, t3.custom_fee
		FROM %s t1
		JOIN %s t2
		ON t1.delivery_id = t2.id
		JOIN %s t3
		ON t1.payment_id = t2.id
	;`, orderTable, deliveryTable, paymentTable)
	rows, err := o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		var delivery model.Delivery
		var payment model.Payment
		if err := rows.Scan(&order.Id, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerId,
			&order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email, &payment.Transaction, &payment.RequestId, &payment.Currency,
			&payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal,
			&payment.CustomFee); err != nil {
			return nil, err
		}
		order.Delivery = delivery
		order.Payment = payment
		orders[order.Id] = order
	}

	query = fmt.Sprintf(`
		SELECT t1.order_uid, t3.chrt_id, t3.track_number, t3.price, t3.rid, t3.item_name, t3.sale, t3.size, t3.total_price, t3.nm_id, t3.brand, t3.status_id
		FROM %s t1
		JOIN %s t2
		ON t1.id = t2.sale_id
		JOIN %s t3
		ON t2.item_id = t3.id
	;`, orderTable, orderItemsTable, itemTable)
	rows, err = o.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	itemMap := make(map[string][]model.Item)
	for rows.Next() {
		var item model.Item
		var orderUid string
		if err := rows.Scan(&orderUid, &item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmId, &item.Brand, &item.Status); err != nil {
			return nil, err
		}
		if _, exist := itemMap[orderUid]; !exist {
			itemMap[orderUid] = make([]model.Item, 0)
		}
		itemMap[orderUid] = append(itemMap[orderUid], item)
	}
	for orderUid, items := range itemMap {
		order, exist := orders[orderUid]
		if exist {
			order.Items = items
			orders[orderUid] = order
		}
	}
	return orders, nil
}

func (o *OrderSql) AddOrder(order model.Order) error {
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (delivery_name, phone, zip, city, delivery_address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID;", deliveryTable)
	row := tx.QueryRow(query, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email)
	var deliveryId int
	if err := row.Scan(&deliveryId); err != nil {
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (transaction_id, request_id, currency, provider_name, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING ID;", paymentTable)
	row = tx.QueryRow(query, order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	var paymentId int
	if err := row.Scan(&paymentId); err != nil {
		return err
	}

	query = fmt.Sprintf(`INSERT INTO %s (order_uid, track_number, entry_name, delivery_id, payment_id, locale, internal_signature, customer_id, 
		delivery_service, shardkey, sm_id, date_created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING ID;`, orderTable)
	row = tx.QueryRow(query, order.Id, order.TrackNumber, order.Entry, deliveryId, paymentId, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated)
	var orderId int
	if err := row.Scan(&orderId); err != nil {
		return err
	}

	for _, item := range order.Items {
		query = fmt.Sprintf(`INSERT INTO %s (chrt_id, track_number, price, rid, item_name, sale, size, total_price, nm_id, brand, status_id) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING ID;`, itemTable)
		row := tx.QueryRow(query, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmId, item.Brand, item.Status)
		var itemId int
		if err := row.Scan(&itemId); err != nil {
			return err
		}

		query = fmt.Sprintf("INSERT INTO %s (sale_id, item_id) VALUES ($1, $2);", orderItemsTable)
		if _, err = tx.Exec(query, orderId, itemId); err != nil {
			return nil
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
