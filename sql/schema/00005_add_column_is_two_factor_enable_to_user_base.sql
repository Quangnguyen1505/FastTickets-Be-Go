-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_base_9999
ADD COLUMN is_two_factor_enabled BOOLEAN DEFAULT FALSE;
COMMENT ON COLUMN pre_go_acc_user_base_9999.is_two_factor_enabled IS 'Authentication is enabled for the user';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_base_9999
DROP COLUMN is_two_factor_enabled;
-- +goose StatementEnd
