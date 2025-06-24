-- name: CreateArticle :one
INSERT INTO articles (title, content, author_id) VALUES ($1, $2, $3) RETURNING id, title, content, author_id, created_at, updated_at;

-- name: GetArticleByID :one
SELECT id, title, content, author_id, created_at, updated_at FROM articles WHERE id = $1 LIMIT 1;

-- name: ListArticles :many
SELECT id, title, content, author_id, created_at, updated_at FROM articles ORDER BY created_at DESC;

-- name: UpdateArticle :one
UPDATE articles SET title = $2, content = $3, updated_at = NOW() WHERE id = $1 RETURNING id, title, content, author_id, created_at, updated_at;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = $1;

-- name: ListArticlesByAuthorID :many
SELECT id, title, content, author_id, created_at, updated_at FROM articles WHERE author_id = $1 ORDER BY created_at DESC;
