-- name: CreateSetting :exec
INSERT INTO settings
    (subsystem, key, value)
 VALUES (
        $1, $2, $3
         ) on conflict (subsystem, key) do update set value = EXCLUDED.value;

-- name: GetSetting :one
select value from settings where subsystem = $1 and key = $2;
