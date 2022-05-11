-- +goose Up
INSERT INTO banners (id, description)
VALUES (1, 'Youth banner'),
       (2, 'Sport banner'),
       (3, 'Health banner');
INSERT INTO slots (id, description)
VALUES (1, 'Top slot'),
       (2, 'Middle slot'),
       (3, 'Bottom banner');
INSERT INTO social_groups (id, description)
VALUES (1, 'Молодежь'),
       (2, 'Люди среднего возраста'),
       (3, 'Пожилые');

-- +goose Down
TRUNCATE TABLE banners;
TRUNCATE TABLE slots;
TRUNCATE TABLE social_groups;
