-- +goose Up
-- +goose StatementBegin
INSERT INTO pre_go_marketing_content (name, description, image, active)
VALUES 
    ('Multi-format support', 'Our platform offers multi-format support, enabling seamless conversions across a wide range of file types for maximum versatility.', 'fm-1.png', true),
    ('Fast and easy', 'Just drop your files on the page, choose an output format and click ''Convert'' button. Wait a little for the process to complete. We aim to do all our conversions in under 1-2 minutes.', 'fm-2.png', true),
    ('In the cloud', 'All conversions take place in the cloud and will not consume any capacity from your computer.', 'fm-3.png', true),
    ('Custom settings', 'Most conversion types support advanced options. For example, with an image converter, you can choose quality, aspect ratio, codec, and other settings, rotate, and flip.', 'fm-4.png', true),
    ('Security guaranteed', 'We delete uploaded files instantly and converted ones after 24 hours. No one has access to your files, and privacy is 100% guaranteed. Read more about security.', 'fm-5.png', true),
    ('All devices supported', 'Converter Free is browser-based and works for all platforms. There is no need to download and install any software.', 'fm-6.png', true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM pre_go_marketing_content 
WHERE name IN (
    'Multi-format support', 
    'Fast and easy', 
    'In the cloud', 
    'Custom settings', 
    'Security guaranteed', 
    'All devices supported'
);
-- +goose StatementEnd
