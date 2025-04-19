-- +goose Up
-- +goose StatementBegin
CREATE TYPE two_factor_auth_type_enum AS ENUM ('SMS', 'EMAIL', 'APP');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_two_factor_9999 (
    two_factor_id SERIAL PRIMARY KEY, -- Khóa chính tự động tăng
    user_id INTEGER NOT NULL, -- Khóa ngoại liên kết tới bảng người dùng
    two_factor_auth_type two_factor_auth_type_enum NOT NULL, -- Loại phương thức 2FA sử dụng ENUM
    two_factor_auth_secret VARCHAR(255) NOT NULL, -- Thông tin bí mật cho 2FA
    two_factor_phone VARCHAR(20), -- Số điện thoại 2FA qua SMS (nếu áp dụng)
    two_factor_email VARCHAR(255), -- Địa chỉ email 2FA qua Email (nếu áp dụng)
    two_factor_is_active BOOLEAN NOT NULL DEFAULT TRUE, -- Trạng thái kích hoạt
    two_factor_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Thời điểm tạo
    two_factor_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Thời điểm cập nhật
    -- Ràng buộc khóa ngoại
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES pre_go_acc_user_base_9999(user_id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_user_id ON pre_go_acc_user_two_factor_9999 (user_id);
CREATE INDEX idx_auth_type ON pre_go_acc_user_two_factor_9999 (two_factor_auth_type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_go_acc_user_two_factor_9999;
DROP TYPE IF EXISTS two_factor_auth_type_enum;
-- +goose StatementEnd
