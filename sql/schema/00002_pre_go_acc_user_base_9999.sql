-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_base_9999  (
    user_id SERIAL PRIMARY KEY,                    -- User ID
    user_account VARCHAR(255) NOT NULL,            -- User account (used to verify identity)
    user_password VARCHAR(255) NOT NULL,           -- User password
    user_salt VARCHAR(255) NOT NULL,               -- Salt used for password encryption
    -- isTwoFactorEnabled
    user_login_time TIMESTAMP NULL DEFAULT NULL,   -- Last login time
    user_logout_time TIMESTAMP NULL DEFAULT NULL,  -- Last logout time
    user_login_ip VARCHAR(45) NULL,                -- Login IP address (45 characters to support IPv6)

    user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Record creation time
    user_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   -- Record update time
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Tạo trigger function để tự động cập nhật cột `user_updated_at` khi bản ghi bị thay đổi
CREATE OR REPLACE FUNCTION update_user_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.user_updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
-- Tạo trigger cho bảng `pre_go_acc_user_base_9999`
CREATE TRIGGER update_user_updated_at 
BEFORE UPDATE ON pre_go_acc_user_base_9999 
FOR EACH ROW 
EXECUTE FUNCTION update_user_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_user_updated_at ON pre_go_acc_user_base_9999;
DROP FUNCTION IF EXISTS update_user_updated_at_column();
DROP TABLE IF EXISTS pre_go_acc_user_base_9999;
-- +goose StatementEnd
