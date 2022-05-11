-- +goose Up
CREATE TABLE IF NOT EXISTS social_groups
(
    id          serial NOT NULL,
    description text   NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE social_groups;