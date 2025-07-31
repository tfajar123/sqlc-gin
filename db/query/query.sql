-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetListUsers :many
SELECT * FROM users ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET name = $2, email = $3 WHERE id = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: CreatePost :one
INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING *;

-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- name: GetListPosts :many
SELECT * FROM posts;

-- name: GetListPostsByUserID :many
SELECT * FROM posts WHERE user_id = $1;

-- name: UpdatePost :one
UPDATE posts SET title = $2, content = $3 WHERE id = $1 RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;