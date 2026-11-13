-- name: CreateSite :one
INSERT INTO bookmarks (
    id, user_id, category_id, title, url, description, logo_url,
    sort_order, is_pinned
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, user_id, category_id, title, url, description, logo_url,
          sort_order, is_pinned, visit_count, last_visited_at,
          created_at, updated_at;

-- name: ListSites :many
SELECT id, user_id, category_id, title, url, description, logo_url,
       sort_order, is_pinned, visit_count, last_visited_at,
       created_at, updated_at
FROM bookmarks
WHERE (sqlc.narg(category_id)::uuid IS NULL OR category_id = sqlc.narg(category_id)::uuid)
  AND (sqlc.narg(user_id)::uuid IS NULL OR user_id = sqlc.narg(user_id)::uuid)
  AND (sqlc.arg(title)::text = '' OR strpos(lower(title), lower(sqlc.arg(title)::text)) > 0)
  AND (sqlc.arg(url)::text = '' OR strpos(lower(url), lower(sqlc.arg(url)::text)) > 0)
  AND (sqlc.arg(description)::text = '' OR strpos(lower(description), lower(sqlc.arg(description)::text)) > 0)
  AND (sqlc.narg(is_pinned)::boolean IS NULL OR is_pinned = sqlc.narg(is_pinned)::boolean)
ORDER BY sort_order ASC, created_at ASC, id ASC
OFFSET sqlc.arg(skip_count)::bigint
LIMIT sqlc.arg(page_limit)::bigint;

-- name: CountSites :one
SELECT count(*) FROM bookmarks
WHERE (sqlc.narg(category_id)::uuid IS NULL OR category_id = sqlc.narg(category_id)::uuid)
  AND (sqlc.narg(user_id)::uuid IS NULL OR user_id = sqlc.narg(user_id)::uuid)
  AND (sqlc.arg(title)::text = '' OR strpos(lower(title), lower(sqlc.arg(title)::text)) > 0)
  AND (sqlc.arg(url)::text = '' OR strpos(lower(url), lower(sqlc.arg(url)::text)) > 0)
  AND (sqlc.arg(description)::text = '' OR strpos(lower(description), lower(sqlc.arg(description)::text)) > 0)
  AND (sqlc.narg(is_pinned)::boolean IS NULL OR is_pinned = sqlc.narg(is_pinned)::boolean)
;

-- name: GetSite :one
SELECT id, user_id, category_id, title, url, description, logo_url,
       sort_order, is_pinned, visit_count, last_visited_at,
       created_at, updated_at
FROM bookmarks
WHERE id = $1
LIMIT 1;

-- name: UpdateSite :one
UPDATE bookmarks
SET
    category_id = CASE WHEN sqlc.arg(set_category_id)::boolean THEN sqlc.arg(category_id)::uuid ELSE category_id END,
    title = CASE WHEN sqlc.arg(set_title)::boolean THEN sqlc.arg(title)::varchar ELSE title END,
    url = CASE WHEN sqlc.arg(set_url)::boolean THEN sqlc.arg(url)::text ELSE url END,
    description = CASE WHEN sqlc.arg(set_description)::boolean THEN sqlc.arg(description)::text ELSE description END,
    logo_url = CASE WHEN sqlc.arg(set_logo_url)::boolean THEN sqlc.arg(logo_url)::text ELSE logo_url END,
    sort_order = CASE WHEN sqlc.arg(set_sort_order)::boolean THEN sqlc.arg(sort_order)::integer ELSE sort_order END,
    is_pinned = CASE WHEN sqlc.arg(set_is_pinned)::boolean THEN sqlc.arg(is_pinned)::boolean ELSE is_pinned END,
    updated_at = now()
WHERE id = sqlc.arg(id)::uuid
RETURNING id, user_id, category_id, title, url, description, logo_url,
          sort_order, is_pinned, visit_count, last_visited_at,
          created_at, updated_at;

-- name: DeleteSite :one
DELETE FROM bookmarks
WHERE id = $1
RETURNING id, user_id, category_id, title, url, description, logo_url,
          sort_order, is_pinned, visit_count, last_visited_at,
          created_at, updated_at;
