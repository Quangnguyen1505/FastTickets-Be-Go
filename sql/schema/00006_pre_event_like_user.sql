-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_event_like_user (
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    liked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (event_id, user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_event_like_user;
-- +goose StatementEnd
