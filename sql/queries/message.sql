-- name: SendMessage :one
INSERT INTO messages (id, send_at, sender_id, receiver_id, content) 
VALUES (
   $1, 
   NOW(),
   $2,
   $3,
   $4
)
RETURNING *;

-- name: GetMessaged :many
SELECT * FROM messages
WHERE sender_id = $1 and receiver_id = $2;