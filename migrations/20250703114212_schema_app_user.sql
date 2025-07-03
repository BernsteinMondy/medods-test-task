-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA app_user;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA app_user;
-- +goose StatementEnd
