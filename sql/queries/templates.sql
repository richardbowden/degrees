-- name: ListTemplates :many
SELECT * from template;

-- name: GetTemplateByRef :one 
SELECT content
FROM template
WHERE ref = $1
  AND deleted_at IS NULL
ORDER BY version DESC
LIMIT 1;

-- name: GetNotificationTemplateByName :one
SELECT t.content
FROM notification_template nt
         JOIN template t on nt.template_id = t.id
WHERE nt.name = $1;

-- name: GetTemplateByID :one
select * from template where id = $1;

-- -- name: ListSystemNotificationTemplates :many
-- select name from notification_template;

-- name: ListSystemNotificationTemplates :many
select t.id, t.name, t.content,t.version, nt.name system_name
    from template t
         join notification_template nt on t.id = nt.id;

-- name: SaveTemplate :exec
insert into template (name, ref, content, scope_type, version)
values ($1, $2, $3, $4, $5);
