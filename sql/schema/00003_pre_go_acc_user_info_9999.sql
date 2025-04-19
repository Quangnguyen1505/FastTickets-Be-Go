-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_info_9999 (
    user_id BIGSERIAL PRIMARY KEY,  -- Primary key for user ID
    user_account VARCHAR(255) NOT NULL,  -- Account of the user
    user_nickname VARCHAR(255),  -- Nickname of the user
    user_avatar VARCHAR(255),  -- Avatar image URL for the user
    user_state SMALLINT NOT NULL DEFAULT 0,  -- User state: 0=Locked, 1=Activated, 2=Not Activated
    user_mobile VARCHAR(20),  -- User's mobile phone number
    
    user_gender SMALLINT,  -- User gender: 0=Secret, 1=Male, 2=Female
    user_birthday DATE,  -- Date of birth
    user_email VARCHAR(255),  -- Email address
    
    user_is_authentication SMALLINT NOT NULL DEFAULT 0,  -- Authentication status: 0=Not Authenticated, 1=Pending, 2=Authenticated

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Record creation time
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Record update time (will be handled by trigger)
);
-- +goose StatementEnd

-- Create function to handle updating 'updated_at' column
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- Create trigger to automatically update 'updated_at' on row update
-- +goose StatementBegin
CREATE TRIGGER update_user_info_updated_at
BEFORE UPDATE ON pre_go_acc_user_info_9999
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_user_info_updated_at ON pre_go_acc_user_info_9999;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP TABLE IF EXISTS pre_go_acc_user_info_9999;
-- +goose StatementEnd
