-- Remove composite unique and restore unique on name
ALTER TABLE teams
    DROP CONSTRAINT IF EXISTS teams_name_sport_id_key;

-- Recreate unique on name (matches original schema)
ALTER TABLE teams
    ADD CONSTRAINT teams_name_key UNIQUE (name);

