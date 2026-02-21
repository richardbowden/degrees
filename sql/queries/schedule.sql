-- name: GetScheduleConfig :many
SELECT * FROM schedule_config
ORDER BY day_of_week;

-- name: GetScheduleConfigForDay :one
SELECT * FROM schedule_config
WHERE day_of_week = $1;

-- name: UpdateScheduleConfig :one
INSERT INTO schedule_config (day_of_week, open_time, close_time, is_open, buffer_minutes)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (day_of_week)
DO UPDATE SET open_time = EXCLUDED.open_time,
              close_time = EXCLUDED.close_time,
              is_open = EXCLUDED.is_open,
              buffer_minutes = EXCLUDED.buffer_minutes
RETURNING *;

-- name: ListBlackoutDates :many
SELECT * FROM schedule_blackouts
ORDER BY date;

-- name: CreateBlackout :one
INSERT INTO schedule_blackouts (date, reason)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteBlackout :exec
DELETE FROM schedule_blackouts
WHERE id = $1;

-- name: IsDateBlackedOut :one
SELECT EXISTS(
    SELECT 1 FROM schedule_blackouts WHERE date = $1
) AS is_blacked_out;
