-- +goose Up

-- --- ENUM TYPES ---

CREATE TYPE order_status AS ENUM ('pending', 'preparing', 'ready', 'delivered', 'cancelled');
CREATE TYPE employee_role AS ENUM ('waiter', 'chef', 'manager');
CREATE TYPE payment_method AS ENUM ('cash', 'card', 'mobile');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded');
CREATE TYPE table_status AS ENUM ('free', 'occupied', 'reserved');

-- --- TABLES ---

CREATE TABLE food_ingredient
(
    id          integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        varchar(255) NOT NULL,
    description text
);

CREATE TABLE allergen
(
    id          integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        varchar(255) NOT NULL,
    description text,
    code        varchar(255) NOT NULL
);

CREATE TABLE restaurant
(
    id      integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name    varchar(255) NOT NULL,
    address text,
    phone   varchar(255),
    email   varchar(255)
);

CREATE TABLE food_items
(
    id             integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id  integer      NOT NULL,
    name           varchar(255) NOT NULL,
    description    text,
    price          numeric(10, 2),
    vat            numeric(5, 2),
    is_available   boolean      NOT NULL DEFAULT true,
    image_url      text,
    prep_time_mins integer
);

CREATE TABLE food_menu
(
    id            integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id integer      NOT NULL,
    name          varchar(255) NOT NULL,
    description   text,
    image_url     text
);

CREATE TABLE employee
(
    id            integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id integer       NOT NULL,
    name          varchar(255),
    role          employee_role NOT NULL,
    pin_hash      varchar(255),
    hire_date     timestamp,
    last_login    timestamp,
    is_active     boolean       NOT NULL DEFAULT true,
    email         varchar(255),
    phone_number  varchar(255)
);

CREATE TABLE food_item_to_menu
(
    menu_id integer NOT NULL,
    item_id integer NOT NULL,
    PRIMARY KEY (menu_id, item_id)
);

CREATE TABLE tables
(
    id       integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    number   integer      NOT NULL,
    capacity integer      NOT NULL,
    section  varchar(255),
    status   table_status NOT NULL DEFAULT 'free'
);

CREATE TABLE food_orders
(
    id                  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id       integer      NOT NULL,
    table_id            integer      NOT NULL,
    waiter_id           integer      NOT NULL,
    order_creation_date timestamp    NOT NULL,
    order_status        order_status NOT NULL,
    order_number        integer      NOT NULL,
    notes               text
);

CREATE TABLE restaurant_ingredient
(
    restaurant_id integer NOT NULL,
    ingredient_id integer NOT NULL,
    is_shared     boolean NOT NULL DEFAULT false,
    PRIMARY KEY (restaurant_id, ingredient_id)
);

CREATE TABLE restaurant_allergen
(
    restaurant_id integer NOT NULL,
    allergen_id   integer NOT NULL,
    PRIMARY KEY (restaurant_id, allergen_id)
);

CREATE TABLE payments
(
    id             integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id       integer        NOT NULL,
    amount         decimal(10, 2) NOT NULL,
    method         payment_method NOT NULL,
    transaction_id varchar(255),
    payment_time   timestamp,
    status         payment_status NOT NULL DEFAULT 'pending'
);

CREATE TABLE ingredient_in_item
(
    ingredient_id integer NOT NULL,
    item_id       integer NOT NULL,
    measurement   text,
    amount        double precision,
    PRIMARY KEY (ingredient_id, item_id)
);

CREATE TABLE item_orders
(
    id                  integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id            integer        NOT NULL,
    item_id             integer        NOT NULL,
    quantity            integer        NOT NULL DEFAULT 1,
    unit_price_snapshot numeric(10, 2) NOT NULL,
    vat_snapshot        numeric(5, 2),
    notes               text
);

CREATE TABLE chat_messages
(
    id            integer PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    restaurant_id integer NOT NULL,
    employee_id   integer NOT NULL,
    message       text    NOT NULL,
    sent_at       timestamp
);

-- --- INDEXES ---

CREATE UNIQUE INDEX employee_index_0 ON employee (id, restaurant_id);

CREATE INDEX item_orders_index_1 ON item_orders (order_id, item_id);

-- --- FOREIGN KEYS ---

ALTER TABLE food_items
    ADD CONSTRAINT food_items_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE food_menu
    ADD CONSTRAINT food_menu_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE employee
    ADD CONSTRAINT employee_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE food_orders
    ADD CONSTRAINT food_orders_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE restaurant_ingredient
    ADD CONSTRAINT restaurant_ingredient_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE restaurant_ingredient
    ADD CONSTRAINT restaurant_ingredient_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES food_ingredient (id) ON DELETE CASCADE;

ALTER TABLE restaurant_allergen
    ADD CONSTRAINT restaurant_allergen_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE restaurant_allergen
    ADD CONSTRAINT restaurant_allergen_allergen_id_fkey FOREIGN KEY (allergen_id) REFERENCES allergen (id) ON DELETE CASCADE;

ALTER TABLE food_orders
    ADD CONSTRAINT food_orders_table_id_fkey FOREIGN KEY (table_id) REFERENCES tables (id) ON DELETE CASCADE;

ALTER TABLE food_orders
    ADD CONSTRAINT food_orders_waiter_id_fkey FOREIGN KEY (waiter_id, restaurant_id) REFERENCES employee (id, restaurant_id) ON DELETE RESTRICT;

ALTER TABLE payments
    ADD CONSTRAINT payments_order_id_fkey FOREIGN KEY (order_id) REFERENCES food_orders (id) ON DELETE CASCADE;

ALTER TABLE item_orders
    ADD CONSTRAINT item_orders_order_id_fkey FOREIGN KEY (order_id) REFERENCES food_orders (id) ON DELETE CASCADE;

ALTER TABLE item_orders
    ADD CONSTRAINT item_orders_item_id_fkey FOREIGN KEY (item_id) REFERENCES food_items (id) ON DELETE CASCADE;

ALTER TABLE food_item_to_menu
    ADD CONSTRAINT food_item_to_menu_item_id_fkey FOREIGN KEY (item_id) REFERENCES food_items (id) ON DELETE CASCADE;

ALTER TABLE food_item_to_menu
    ADD CONSTRAINT food_item_to_menu_menu_id_fkey FOREIGN KEY (menu_id) REFERENCES food_menu (id) ON DELETE CASCADE;

ALTER TABLE ingredient_in_item
    ADD CONSTRAINT ingredient_in_item_ingredient_id_fkey FOREIGN KEY (ingredient_id) REFERENCES food_ingredient (id) ON DELETE CASCADE;

ALTER TABLE ingredient_in_item
    ADD CONSTRAINT ingredient_in_item_item_id_fkey FOREIGN KEY (item_id) REFERENCES food_items (id) ON DELETE CASCADE;

ALTER TABLE chat_messages
    ADD CONSTRAINT chat_messages_restaurant_id_fkey FOREIGN KEY (restaurant_id) REFERENCES restaurant (id) ON DELETE CASCADE;

ALTER TABLE chat_messages
    ADD CONSTRAINT chat_messages_employee_id_fkey FOREIGN KEY (employee_id, restaurant_id) REFERENCES employee (id, restaurant_id) ON DELETE CASCADE;

-- +goose Down

DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS item_orders CASCADE;
DROP TABLE IF EXISTS ingredient_in_item CASCADE;
DROP TABLE IF EXISTS payments CASCADE;
DROP TABLE IF EXISTS restaurant_allergen CASCADE;
DROP TABLE IF EXISTS restaurant_ingredient CASCADE;
DROP TABLE IF EXISTS food_orders CASCADE;
DROP TABLE IF EXISTS tables CASCADE;
DROP TABLE IF EXISTS food_item_to_menu CASCADE;
DROP TABLE IF EXISTS employee CASCADE;
DROP TABLE IF EXISTS food_menu CASCADE;
DROP TABLE IF EXISTS food_items CASCADE;
DROP TABLE IF EXISTS restaurant CASCADE;
DROP TABLE IF EXISTS allergen CASCADE;
DROP TABLE IF EXISTS food_ingredient CASCADE;

DROP TYPE IF EXISTS table_status;
DROP TYPE IF EXISTS payment_status;
DROP TYPE IF EXISTS payment_method;
DROP TYPE IF EXISTS employee_role;
DROP TYPE IF EXISTS order_status;