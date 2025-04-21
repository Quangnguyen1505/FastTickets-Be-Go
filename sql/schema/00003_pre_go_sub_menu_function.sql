-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_sub_menu_function (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image TEXT,
    url TEXT,
    menufuncId VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_menufunc FOREIGN KEY (menufuncId) REFERENCES pre_go_menu_function(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_go_sub_menu_function;
-- +goose StatementEnd