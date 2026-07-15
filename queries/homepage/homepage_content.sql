-- name: GetHomepageContent :one
SELECT * FROM homepage_content
WHERE id = 1;

-- name: GetHomepageContentForDisplay :one
SELECT
    CASE WHEN sqlc.arg(locale)::text = 'alt' THEN COALESCE(welcome_text_alt, welcome_text) ELSE welcome_text END::text AS welcome_text,
    image_one_url,
    image_two_url
FROM homepage_content
WHERE id = 1;

-- name: UpdateHomepageContent :exec
UPDATE homepage_content
SET welcome_text = $1, welcome_text_alt = $2, image_one_url = $3, image_two_url = $4, updated_at = now()
WHERE id = 1;