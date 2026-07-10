-- +goose Up

ALTER TABLE food_items ADD COLUMN instructions text;

CREATE TABLE food_item_allergen (
                                    item_id integer NOT NULL,
                                    allergen_id integer NOT NULL,
                                    PRIMARY KEY (item_id, allergen_id)
);

ALTER TABLE food_item_allergen ADD CONSTRAINT food_item_allergen_item_id_fkey
    FOREIGN KEY (item_id) REFERENCES food_items (id) ON DELETE CASCADE;

ALTER TABLE food_item_allergen ADD CONSTRAINT food_item_allergen_allergen_id_fkey
    FOREIGN KEY (allergen_id) REFERENCES allergen (id) ON DELETE CASCADE;

-- +goose Down

DROP TABLE IF EXISTS food_item_allergen;
ALTER TABLE food_items DROP COLUMN IF EXISTS instructions;