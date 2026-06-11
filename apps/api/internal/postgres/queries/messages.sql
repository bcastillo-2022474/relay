-- name: InsertMessage :exec
INSERT INTO messages (id, organization_id, application_id, event_type_id, payload, status)
VALUES (@id, @organization_id, @application_id, @event_type_id, @payload, @status);

-- name: ClaimMessageBatch :many
UPDATE messages
SET attempts = attempts + 1
WHERE id IN (
    SELECT id FROM messages
    WHERE status = 'pending' AND next_attempt_at <= now()
    ORDER BY created_at
    LIMIT @batch_size
    FOR UPDATE SKIP LOCKED
)
RETURNING id, organization_id, application_id, event_type_id, payload, status;

-- name: MarkMessagePublished :exec
UPDATE messages
SET status = 'published', published_at = now()
WHERE id = @id;

-- name: RescheduleMessage :exec
UPDATE messages
SET next_attempt_at = @next_attempt_at
WHERE id = @id;

-- name: MarkMessageFailed :exec
UPDATE messages
SET status = 'failed'
WHERE id = @id;
