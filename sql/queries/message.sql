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