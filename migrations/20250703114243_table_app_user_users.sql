-- +goose Up
-- +goose StatementBegin
CREATE TABLE app_user.users
(
    id              uuid NOT NULL PRIMARY KEY,
    username        text NOT NULL,
    hash            text NOT NULL,
    refresh_token   text NOT NULL,
    access_token_id uuid NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE app_user.users;
-- +goose StatementEnd
