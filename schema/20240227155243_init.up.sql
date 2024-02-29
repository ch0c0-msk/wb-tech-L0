CREATE TABLE DELIVERY (
    id SERIAL PRIMARY KEY,
    delivery_name VARCHAR(64) NOT NULL,
    phone VARCHAR(16) NOT NULL,
    zip VARCHAR(16) NOT NULL,
    city VARCHAR(32) NOT NULL,
    delivery_address VARCHAR(64) NOT NULL,
    region VARCHAR(32) NOT NULL,
    email VARCHAR(64) NOT  NULL
);

CREATE TABLE PAYMENT (
    id SERIAL PRIMARY KEY,
    transaction_id VARCHAR(64) NOT NULL,
    request_id VARCHAR(64),
    currency VARCHAR(8) NOT NULL,
    provider_name VARCHAR(64) NOT NULL,
    amount INT NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(64) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL
);

CREATE TABLE SALE (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(32) UNIQUE NOT NULL,
    track_number VARCHAR(32) NOT NULL,
    entry_name VARCHAR(16) NOT NULL,
    delivery_id INT NOT NULL,
    payment_id INT NOT NULL,
    locale VARCHAR(4) NOT NULL,
    internal_signature VARCHAR(16),
    customer_id VARCHAR(32) NOT NULL,
    delivery_service VARCHAR(64) NOT NULL,
    shardkey VARCHAR(4) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (delivery_id) REFERENCES delivery(id),
    FOREIGN KEY (payment_id) REFERENCES payment(id)
);

CREATE TABLE ITEM (
    id SERIAL PRIMARY KEY,
    chrt_id INT NOT NULL,
    track_number VARCHAR(32) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(64) NOT NULL,
    item_name VARCHAR(64) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(8) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(64) NOT NULL,
    status_id INT NOT NULL
);

CREATE TABLE SALE_ITEMS (
    id SERIAL PRIMARY KEY,
    sale_id INT NOT NULL,
    item_id INT NOT NULL,
    FOREIGN KEY (sale_id) REFERENCES sale(id),
    FOREIGN KEY (item_id) REFERENCES item(id)
);

