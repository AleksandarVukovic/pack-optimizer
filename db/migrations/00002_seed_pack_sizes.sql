-- +goose Up
INSERT INTO pack_sizes (size) VALUES (250), (500), (1000), (2000), (5000);

-- +goose Down
DELETE FROM pack_sizes WHERE size IN (250, 500, 1000, 2000, 5000);
