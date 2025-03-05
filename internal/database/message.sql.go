// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: message.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getMessaged = `-- name: GetMessaged :many
SELECT id, sent_at, sender_id, receiver_id, content FROM messages
WHERE sender_id = $1 and receiver_id = $2
`

type GetMessagedParams struct {
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
}

func (q *Queries) GetMessaged(ctx context.Context, arg GetMessagedParams) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, getMessaged, arg.SenderID, arg.ReceiverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.SentAt,
			&i.SenderID,
			&i.ReceiverID,
			&i.Content,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const sendMessage = `-- name: SendMessage :one
INSERT INTO messages (id, send_at, sender_id, receiver_id, content) 
VALUES (
   $1, 
   NOW(),
   $2,
   $3,
   $4
)
RETURNING id, sent_at, sender_id, receiver_id, content
`

type SendMessageParams struct {
	ID         uuid.UUID
	SenderID   uuid.UUID
	ReceiverID uuid.UUID
	Content    string
}

func (q *Queries) SendMessage(ctx context.Context, arg SendMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, sendMessage,
		arg.ID,
		arg.SenderID,
		arg.ReceiverID,
		arg.Content,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.SentAt,
		&i.SenderID,
		&i.ReceiverID,
		&i.Content,
	)
	return i, err
}
