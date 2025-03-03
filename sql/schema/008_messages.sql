-- +goose Up
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    sent_at TIMESTAMP NOT NULL,
    sender_id VARCHAR(255) NOT NULL,
    receiver_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL
)

-- +goose Down
DROP TABLE messages;