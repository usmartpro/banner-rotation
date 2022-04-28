-- +goose Up
CREATE TABLE IF NOT EXISTS slots
(
    id          serial NOT NULL,
    description text   NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE slots;