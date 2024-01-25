-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TYPE type AS ENUM ('CREDENTIALS', 'TEXT', 'BINARY', 'CARD');
CREATE TABLE IF NOT EXISTS data (
    id serial PRIMARY KEY,
    user_id integer REFERENCES users (id),
    name varchar(64),
    type type,
    data bytea,
    created_at timestamp DEFAULT NOW(),
    metadata json
);

-- create indexes
CREATE INDEX IF NOT EXISTS data_user_id_idx ON data (user_id);
CREATE INDEX IF NOT EXISTS data_name_idx ON data (name);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP INDEX data_name_idx;
DROP INDEX data_user_id_idx;
DROP TABLE data;
DROP TYPE type;
