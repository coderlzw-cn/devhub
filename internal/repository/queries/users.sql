-- name: CreateUser :one
INSERT INTO users (username, email, age, password)
VALUES ($1, $2, $3, $4)
RETURNING id, username, email, age, password;

-- name: CreateUserWithID :one
INSERT INTO users (id, username, email, age, password)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, username, email, age, password;

-- name: ListUsers :many
SELECT id, username, email, age, password
FROM users
ORDER BY id ASC
OFFSET sqlc.arg(skip_count)::bigint
LIMIT sqlc.arg(page_limit)::bigint;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: GetUserByUsername :one
SELECT id, username, email, age, password
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUser :one
SELECT id, username, email, age, password
FROM users
WHERE id = $1
LIMIT 1;

-- name: UpdateUser :execrows
UPDATE users
SET
    username = CASE WHEN sqlc.arg(set_username)::boolean THEN sqlc.arg(username)::varchar ELSE username END,
    email = CASE WHEN sqlc.arg(set_email)::boolean THEN sqlc.arg(email)::varchar ELSE email END,
    age = CASE WHEN sqlc.arg(set_age)::boolean THEN sqlc.arg(age)::integer ELSE age END
WHERE id = sqlc.arg(id)::bigint;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id, username, email, age, password;
