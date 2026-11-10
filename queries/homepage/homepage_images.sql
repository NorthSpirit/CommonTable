-- name: ListHomepageImages :many
SELECT * FROM homepage_images
ORDER BY slot;

-- name: UpsertHomepageImage :exec
INSERT INTO homepage_images (slot, url, width_percent, shape, position_percent)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (slot) DO UPDATE
SET url = EXCLUDED.url,
width_percent = EXCLUDED.width_percent,
shape = EXCLUDED.shape,
position_percent = EXCLUDED.position_percent,
updated_at = now();

-- name: DeleteHomepageImage :exec
DELETE FROM homepage_images
WHERE slot = $1;
