CREATE TABLE IF NOT EXISTS notification_queue
(
    id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id           BIGINT       NOT NULL,
    notification_type VARCHAR(50)  NOT NULL,
    template_id       VARCHAR(100) NOT NULL,
    payload           JSONB,
    status            VARCHAR(20)  NOT NULL    DEFAULT 'pending',
    created_at        TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    scheduled_for     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    processed_at      TIMESTAMP WITH TIME ZONE,
    retry_count       INTEGER                  DEFAULT 0,

    CONSTRAINT notification_queue_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed'))
);

CREATE INDEX idx_notification_queue_user_id ON notification_queue (user_id);
CREATE INDEX idx_notification_queue_status ON notification_queue (status);
CREATE INDEX idx_notification_queue_scheduled_for ON notification_queue (scheduled_for);

CREATE TABLE IF NOT EXISTS notification_history
(
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    notification_id BIGINT       NOT NULL,
    status          VARCHAR(20)  NOT NULL,
    changed_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    error_message   TEXT,
    metadata        JSONB,

    -- INSERT, UPDATE, DELETE
    operation_type  VARCHAR(10)  NOT NULL DEFAULT 'UPDATE',

    -- status_change, retry_attempt, reschedule, etc.
    change_type     VARCHAR(20),

    CONSTRAINT notification_history_status_check
        CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    CONSTRAINT notification_history_operation_check
        CHECK (operation_type IN ('INSERT', 'UPDATE', 'DELETE'))
);

CREATE INDEX idx_notification_history_notification_id ON notification_history (notification_id);
CREATE INDEX idx_notification_history_changed_at ON notification_history (changed_at);
CREATE INDEX idx_notification_history_operation_type ON notification_history (operation_type);
CREATE INDEX idx_notification_history_change_type ON notification_history (change_type);

CREATE OR REPLACE FUNCTION log_notification_changes()
RETURNS TRIGGER AS $$
DECLARE
    old_record JSONB;
    new_record JSONB;
    changes JSONB := '{}'::JSONB;
    field_name TEXT;
    change_type TEXT;
BEGIN
    -- Handle INSERT operations
    IF TG_OP = 'INSERT' THEN
        INSERT INTO notification_history (
            notification_id,
            status,
            changed_at,
            metadata,
            operation_type
        ) VALUES (
            NEW.id,
            NEW.status,
            CURRENT_TIMESTAMP,
            jsonb_build_object(
                'operation', 'INSERT',
                'created_record', row_to_json(NEW)::JSONB
            ),
            TG_OP
        );
        RETURN NEW;
    END IF;
    
    -- Handle UPDATE operations
    IF TG_OP = 'UPDATE' THEN
        old_record := row_to_json(OLD)::JSONB;
        new_record := row_to_json(NEW)::JSONB;

        FOR field_name IN SELECT * FROM jsonb_object_keys(new_record) LOOP
            IF old_record->field_name IS DISTINCT FROM new_record->field_name THEN
                changes := changes || jsonb_build_object(
                    field_name, jsonb_build_object(
                        'old_value', old_record->field_name,
                        'new_value', new_record->field_name
                    )
                );
            END IF;
        END LOOP;

        IF changes != '{}'::JSONB THEN
            -- Determine the change
            change_type := CASE
                WHEN changes ? 'status' THEN 'status_change'
                WHEN changes ? 'retry_count' THEN 'retry_attempt'
                WHEN changes ? 'scheduled_for' THEN 'reschedule'
                WHEN changes ? 'processed_at' THEN 'processing_update'
                ELSE 'field_update'
            END;
            
            INSERT INTO notification_history (
                notification_id,
                status,
                changed_at,
                error_message,
                metadata,
                operation_type,
                change_type
            ) VALUES (
                NEW.id,
                NEW.status,
                CURRENT_TIMESTAMP,
                CASE 
                    WHEN NEW.status = 'failed' AND changes ? 'status' THEN 
                        'Status changed to failed'
                    ELSE NULL 
                END,
                jsonb_build_object(
                    'operation', 'UPDATE',
                    'change_type', change_type,
                    'changes', changes,
                    'change_count', jsonb_object_keys(changes)
                ),
                TG_OP,
                change_type
            );
        END IF;
        
        RETURN NEW;
    END IF;
    
    -- DELETE
    IF TG_OP = 'DELETE' THEN
        INSERT INTO notification_history (
            notification_id,
            status,
            changed_at,
            metadata,
            operation_type
        ) VALUES (
            OLD.id,
            OLD.status,
            CURRENT_TIMESTAMP,
            jsonb_build_object(
                'operation', 'DELETE',
                'deleted_record', row_to_json(OLD)::JSONB
            ),
            TG_OP
        );
        RETURN OLD;
    END IF;
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_notification_changes
    AFTER INSERT OR UPDATE OR DELETE ON notification_queue
    FOR EACH ROW
    EXECUTE FUNCTION log_notification_changes();