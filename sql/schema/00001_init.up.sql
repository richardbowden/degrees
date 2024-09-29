CREATE EXTENSION if not exists "uuid-ossp";

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_updated_at_trigger(target_table regclass)
RETURNS VOID AS $$
#variable_conflict use_column
DECLARE
    column_exists BOOLEAN;
    schema_name TEXT;
    table_name TEXT;

BEGIN
    -- Extract schema and table name
    SELECT n.nspname, c.relname INTO schema_name, table_name
    FROM pg_class c
    JOIN pg_namespace n ON n.oid = c.relnamespace
    WHERE c.oid = target_table;

    -- Check if updated_at column exists
    SELECT EXISTS (
        SELECT FROM information_schema.columns
        WHERE table_schema = schema_name
          AND table_name = table_name
          AND column_name = 'updated_at'
    ) INTO column_exists;

    -- If updated_at column doesn't exist, add it
    IF NOT column_exists THEN
        RAISE EXCEPTION 'trying to add update trigger to %I in %I', table_name, schema_name;
    END IF;

    -- Create the trigger
    EXECUTE format('
        CREATE TRIGGER set_updated_at
        BEFORE UPDATE ON %I.%I
        FOR EACH ROW
        EXECUTE FUNCTION update_modified_column();
    ', schema_name, table_name);
END;
$$ LANGUAGE plpgsql;
