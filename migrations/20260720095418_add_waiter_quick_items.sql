-- +goose Up

CREATE TABLE waiter_quick_items (
                                    employee_id integer NOT NULL,
                                    item_id integer NOT NULL,
                                    sort_order integer NOT NULL DEFAULT 0,

                                    PRIMARY KEY (employee_id, item_id)
);

ALTER TABLE waiter_quick_items ADD CONSTRAINT waiter_quick_items_employee_id_fkey
    FOREIGN KEY (employee_id) REFERENCES employee (id) ON DELETE CASCADE;

ALTER TABLE waiter_quick_items ADD CONSTRAINT waiter_quick_items_item_id_fkey
    FOREIGN KEY (item_id) REFERENCES food_items (id) ON DELETE CASCADE;


-- +goose Down

DROP TABLE IF EXISTS waiter_quick_items;