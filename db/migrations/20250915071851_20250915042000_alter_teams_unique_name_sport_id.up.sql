-- Drop existing unique constraint on name if it exists
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM   pg_constraint c
        JOIN   pg_class t ON t.oid = c.conrelid
        WHERE  t.relname = 'teams'
        AND    c.conname = 'teams_name_key'
    ) THEN
        ALTER TABLE teams DROP CONSTRAINT teams_name_key;
    END IF;
END $$;

-- Add composite unique constraint on (name, sport_id)
ALTER TABLE teams
    ADD CONSTRAINT teams_name_sport_id_key UNIQUE (name, sport_id);

