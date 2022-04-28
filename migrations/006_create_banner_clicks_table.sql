-- +goose Up
CREATE TABLE IF NOT EXISTS banner_clicks
(
    id              serial    NOT NULL,
    banner_id       serial    NOT NULL,
    slot_id         serial    NOT NULL,
    social_group_id serial    NOT NULL,
    date            timestamp NOT NULL DEFAULT current_timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY (banner_id)
        REFERENCES banners (id),
    FOREIGN KEY (slot_id)
        REFERENCES slots (id),
    FOREIGN KEY (social_group_id)
        REFERENCES social_groups (id)
);

-- +goose Down
DROP TABLE banner_clicks;