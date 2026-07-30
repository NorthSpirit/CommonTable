-- +goose Up

CREATE UNIQUE INDEX employee_restaurant_name_unique_active
    ON employee (restaurant_id, name)
    WHERE is_active = true;

-- +goose Down

DROP INDEX IF EXISTS employee_restaurant_name_unique_active;