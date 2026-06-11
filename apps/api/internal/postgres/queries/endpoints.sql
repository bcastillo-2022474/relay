-- name: InsertEndpoint :exec
INSERT INTO endpoints (id, application_id, organization_id, url, description, signing_secret)
VALUES (@id, @application_id, @organization_id, @url, @description, @signing_secret);

-- name: FindEndpointByID :one
SELECT *
FROM endpoints
WHERE id = @id AND organization_id = @organization_id
LIMIT 1;

-- name: FindEndpointByURL :one
SELECT *
FROM endpoints
WHERE url = @url AND organization_id = @organization_id
LIMIT 1;

-- name: DeleteEndpoint :exec
DELETE FROM endpoints
WHERE id = @id AND organization_id = @organization_id;

-- name: InsertEndpointSubscription :exec
INSERT INTO endpoint_subscriptions (endpoint_id, event_type_id, organization_id)
VALUES (@endpoint_id, @event_type_id, @organization_id);

-- name: DeleteEndpointSubscription :exec
DELETE FROM endpoint_subscriptions
WHERE endpoint_id = @endpoint_id AND event_type_id = @event_type_id AND organization_id = @organization_id;

-- name: FindSubscribedEventTypes :many
SELECT event_type_id
FROM endpoint_subscriptions
WHERE endpoint_id = @endpoint_id AND organization_id = @organization_id;
