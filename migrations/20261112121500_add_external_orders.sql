-- +goose Up

CREATE TABLE restaurant_api_keys
(
    id            integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id integer      NOT NULL REFERENCES restaurant (id) ON DELETE CASCADE,
    key_hash      varchar(64)  NOT NULL UNIQUE,
    key_prefix    varchar(12)  NOT NULL,
    label         varchar(255) NOT NULL,
    created_at    timestamp    NOT NULL DEFAULT now(),
    revoked_at    timestamp
);

CREATE INDEX restaurant_api_keys_restaurant_idx ON restaurant_api_keys (restaurant_id);

CREATE TABLE unapproved_orders
(
    id                integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id     integer      NOT NULL REFERENCES restaurant (id) ON DELETE CASCADE,
    api_key_id        integer      NOT NULL REFERENCES restaurant_api_keys (id) ON DELETE CASCADE,
    customer_name     varchar(255) NOT NULL,
    customer_phone    varchar(32)  NOT NULL,
    notes             text,
    status            text         NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    idempotency_key   varchar(64)  NOT NULL,
    approved_order_id integer      REFERENCES food_orders (id) ON DELETE SET NULL,
    decided_by        integer      REFERENCES employee (id) ON DELETE SET NULL,
    decided_at        timestamp,
    rejection_reason  text,
    created_at        timestamp    NOT NULL DEFAULT now(),
    UNIQUE (api_key_id, idempotency_key)
);

CREATE INDEX unapproved_orders_key_created_idx ON unapproved_orders (api_key_id, created_at);
CREATE INDEX unapproved_orders_key_status_idx ON unapproved_orders (api_key_id, status);
CREATE INDEX unapproved_orders_phone_status_idx ON unapproved_orders (customer_phone, status);
CREATE INDEX unapproved_orders_restaurant_status_idx ON unapproved_orders (restaurant_id, status, created_at);

CREATE TABLE unapproved_order_items
(
    id                  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    unapproved_order_id integer        NOT NULL REFERENCES unapproved_orders (id) ON DELETE CASCADE,
    item_id             integer        NOT NULL REFERENCES food_items (id) ON DELETE CASCADE,
    quantity            integer        NOT NULL CHECK (quantity BETWEEN 1 AND 1000),
    unit_price_snapshot numeric(10, 2) NOT NULL,
    vat_snapshot        numeric(5, 2),
    item_name_snapshot  varchar(255)   NOT NULL,
    notes               text
);

CREATE INDEX unapproved_order_items_order_idx ON unapproved_order_items (unapproved_order_id);

-- +goose Down

DROP TABLE IF EXISTS unapproved_order_items;
DROP TABLE IF EXISTS unapproved_orders;
DROP TABLE IF EXISTS restaurant_api_keys;
