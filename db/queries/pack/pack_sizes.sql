-- name: ListPackSizes :many
SELECT size FROM pack_sizes ORDER BY size;

-- name: DeletePackSizes :exec
DELETE FROM pack_sizes;

-- name: InsertPackSizes :exec
INSERT INTO pack_sizes (size)
SELECT unnest(sqlc.arg(sizes)::smallint[]);
