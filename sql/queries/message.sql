-- name: SendMessage :one
INSERT INTO messages (id, sent_at, sender_id, receiver_id, content) 
VALUES (
   $1, 
   NOW(),
   $2,
   $3,
   $4
)
RETURNING *;

-- name: GetMessages :many
SELECT * FROM messages
WHERE sender_id = $1 and receiver_id = $2
ORDER BY sent_at;

-- name: DeleteMessage :exec
DELETE FROM messages
WHERE id = $1;