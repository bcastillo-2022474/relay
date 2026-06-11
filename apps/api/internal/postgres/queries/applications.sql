-- name: InsertApplication :exec
INSERT INTO applications (id, name, slug, organization_id)
VALUES (@id, @name, @slug, @organization_id);

-- name: FindApplicationByID :one
SELECT *
FROM applications
WHERE id = @id AND organization_id = @organization_id
LIMIT 1;

-- name: FindApplicationBySlug :one
SELECT *
FROM applications
WHERE slug = @slug AND organization_id = @organization_id
LIMIT 1;

-- name: DeleteApplication :exec
DELETE FROM applications
WHERE id = @id AND organization_id = @organization_id;
