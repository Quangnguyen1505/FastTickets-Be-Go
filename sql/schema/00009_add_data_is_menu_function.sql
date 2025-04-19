-- +goose Up
-- +goose StatementBegin
INSERT INTO pre_go_menu_function (id, name, description, url, active)
VALUES 
    ('menu1', 'All tools', 'kaka1', 'url sample', true),
    ('menu2', 'Convert PDF', 'kaka2', 'url sample', true),
    ('menu3', 'Convert Docx', 'kaka3', 'url sample', true),
    ('menu4', 'Convert Image', 'kaka4', 'url sample', true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM pre_go_menu_function 
WHERE id IN (
    'menu1', 
    'menu2', 
    'menu3', 
    'menu4'
);
-- +goose StatementEnd
