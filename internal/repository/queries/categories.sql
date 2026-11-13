-- name: CreateCategory :one
INSERT INTO categories (
    id, user_id, parent_id, name, slug, description, icon, sort_order, is_visible
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, parent_id, name, slug, description, icon,
          sort_order, is_visible, created_at, updated_at;

-- name: ListCategories :many
SELECT id, user_id, parent_id, name, slug, description, icon,
       sort_order, is_visible, created_at, updated_at
FROM categories
WHERE strpos(lower(name), lower(sqlc.arg(name)::text)) > 0
  AND strpos(lower(slug), lower(sqlc.arg(slug)::text)) > 0
  AND (sqlc.arg(description)::text = '' OR strpos(lower(description), lower(sqlc.arg(description)::text)) > 0)
ORDER BY sort_order ASC, created_at ASC
OFFSET sqlc.arg(skip_count)::bigint
LIMIT sqlc.arg(page_limit)::bigint;

-- name: CountCategories :one
SELECT count(*) FROM categories
WHERE strpos(lower(name), lower(sqlc.arg(name)::text)) > 0
  AND strpos(lower(slug), lower(sqlc.arg(slug)::text)) > 0
  AND (sqlc.arg(description)::text = '' OR strpos(lower(description), lower(sqlc.arg(description)::text)) > 0);

-- name: GetCategory :one
SELECT id, user_id, parent_id, name, slug, description, icon,
       sort_order, is_visible, created_at, updated_at
FROM categories
WHERE id = $1
LIMIT 1;

-- name: UpdateCategory :one
UPDATE categories
SET
    parent_id = CASE WHEN sqlc.arg(set_parent_id)::boolean THEN sqlc.arg(parent_id)::uuid ELSE parent_id END,
    name = CASE WHEN sqlc.arg(set_name)::boolean THEN sqlc.arg(name)::varchar ELSE name END,
    slug = CASE WHEN sqlc.arg(set_slug)::boolean THEN sqlc.arg(slug)::varchar ELSE slug END,
    description = CASE WHEN sqlc.arg(set_description)::boolean THEN sqlc.arg(description)::text ELSE description END,
    icon = CASE WHEN sqlc.arg(set_icon)::boolean THEN sqlc.arg(icon)::varchar ELSE icon END,
    sort_order = CASE WHEN sqlc.arg(set_sort_order)::boolean THEN sqlc.arg(sort_order)::integer ELSE sort_order END,
    is_visible = CASE WHEN sqlc.arg(set_is_visible)::boolean THEN sqlc.arg(is_visible)::boolean ELSE is_visible END,
    updated_at = now()
WHERE id = sqlc.arg(id)::uuid
RETURNING id, user_id, parent_id, name, slug, description, icon,
          sort_order, is_visible, created_at, updated_at;

-- name: DeleteCategory :one
DELETE FROM categories
WHERE id = $1
RETURNING id, user_id, parent_id, name, slug, description, icon,
          sort_order, is_visible, created_at, updated_at;

-- name: CountChildCategories :one
SELECT count(*) FROM categories
WHERE parent_id = $1;

-- name: CountCategorySites :one
SELECT count(*) FROM bookmarks
WHERE category_id = $1;
