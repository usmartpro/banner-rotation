-- +goose Up
CREATE TABLE IF NOT EXISTS banners
(
    id          serial NOT NULL,
    description text   NOT NULL DEFAULT '',
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE banners;