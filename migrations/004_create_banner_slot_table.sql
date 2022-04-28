-- +goose Up
CREATE TABLE IF NOT EXISTS banner_slot
(
    banner_id serial NOT NULL,
    slot_id   serial NOT NULL,
    PRIMARY KEY (banner_id, slot_id),
    FOREIGN KEY (banner_id)
        REFERENCES banners (id),
    FOREIGN KEY (slot_id)
        REFERENCES slots (id)
);

-- +goose Down
DROP TABLE banner_slot;