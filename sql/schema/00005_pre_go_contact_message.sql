-- +goose Up
-- +goose StatementBegin
CREATE TABLE pre_go_contact_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(255),
    message TEXT NOT NULL,
    status SMALLINT NOT NULL DEFAULT 0, -- 0: pending, 1: read, 2: responded 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_go_contact_messages;
-- +goose StatementEnd
