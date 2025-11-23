-- name: CreateSetting :exec
INSERT INTO settings
    (key, value)
 VALUES (
        $1, $2
         ) on conflict (key) do update set value = EXCLUDED.value;

-- name: GetSetting :one
select value from settings where key = $1;

