-- name: CreateSetting :exec
INSERT INTO setting
    (key, value)
 VALUES (
        $1, $2
         ) on conflict (key) do update set value = EXCLUDED.value;

-- name: GetSetting :one
select value from setting where key = $1;

