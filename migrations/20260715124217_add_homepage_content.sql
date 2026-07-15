-- +goose Up

CREATE TABLE homepage_content (
                                  id integer PRIMARY KEY,
                                  welcome_text text NOT NULL,
                                  welcome_text_alt text,
                                  image_one_url text,
                                  image_two_url text,
                                  updated_at timestamp NOT NULL DEFAULT now()
);

-- Seed the single row with sensible defaults so the homepage never renders empty.
INSERT INTO homepage_content (id, welcome_text)
VALUES (1, 'Welcome to CommonTable — browse restaurants, explore their menus, and check what''s in each dish before you order.');

-- +goose Down

DROP TABLE IF EXISTS homepage_content;