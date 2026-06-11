-- name: InsertEventType :exec
INSERT INTO event_types (id, name, application_id, organization_id, payload_schema)
VALUES (@id, @name, @application_id, @organization_id, @payload_schema);

-- name: FindEventTypeByID :one
SELECT *
FROM event_types
WHERE id = @id AND organization_id = @organization_id
LIMIT 1;

-- name: FindEventTypeByName :one
SELECT *
FROM event_types
WHERE name = @name AND application_id = @application_id AND organization_id = @organization_id
LIMIT 1;

-- name: DeleteEventType :exec
DELETE FROM event_types
WHERE id = @id AND organization_id = @organization_id;
