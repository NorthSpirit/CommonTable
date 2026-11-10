-- +goose Up

CREATE TABLE homepage_images (
    slot             integer   PRIMARY KEY CHECK (slot BETWEEN 1 AND 4),
    url              text      NOT NULL,
    width_percent    integer   NOT NULL DEFAULT 50 CHECK (width_percent BETWEEN 10 AND 100),
    shape            text      NOT NULL DEFAULT 'original' CHECK (shape IN ('original', '16:9', '4:3', '1:1', '3:4')),
    position_percent integer   NOT NULL DEFAULT 50 CHECK (position_percent BETWEEN 0 AND 100),
    updated_at       timestamp NOT NULL DEFAULT now()
);

INSERT INTO homepage_images (slot, url)
SELECT 1, image_one_url FROM homepage_content WHERE id = 1 AND COALESCE(image_one_url, '') <> '';

INSERT INTO homepage_images (slot, url)
SELECT 2, image_two_url FROM homepage_content WHERE id = 1 AND COALESCE(image_two_url, '') <> '';

ALTER TABLE homepage_content
    DROP COLUMN image_one_url,
    DROP COLUMN image_two_url;

-- +goose Down

ALTER TABLE homepage_content
    ADD COLUMN image_one_url text,
    ADD COLUMN image_two_url text;

UPDATE homepage_content SET image_one_url = (SELECT url FROM homepage_images WHERE slot = 1) WHERE id = 1;
UPDATE homepage_content SET image_two_url = (SELECT url FROM homepage_images WHERE slot = 2) WHERE id = 1;

DROP TABLE IF EXISTS homepage_images;
