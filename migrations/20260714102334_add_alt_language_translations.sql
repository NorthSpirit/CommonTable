-- +goose Up

ALTER TABLE food_items ADD COLUMN name_alt varchar(255);
ALTER TABLE food_items ADD COLUMN description_alt text;
ALTER TABLE food_items ADD COLUMN instructions_alt text;

ALTER TABLE food_menu ADD COLUMN name_alt varchar(255);
ALTER TABLE food_menu ADD COLUMN description_alt text;

ALTER TABLE allergen ADD COLUMN name_alt varchar(255);
ALTER TABLE allergen ADD COLUMN description_alt text;

-- +goose Down

ALTER TABLE allergen DROP COLUMN IF EXISTS description_alt;
ALTER TABLE allergen DROP COLUMN IF EXISTS name_alt;

ALTER TABLE food_menu DROP COLUMN IF EXISTS description_alt;
ALTER TABLE food_menu DROP COLUMN IF EXISTS name_alt;

ALTER TABLE food_items DROP COLUMN IF EXISTS instructions_alt;
ALTER TABLE food_items DROP COLUMN IF EXISTS description_alt;
ALTER TABLE food_items DROP COLUMN IF EXISTS name_alt;