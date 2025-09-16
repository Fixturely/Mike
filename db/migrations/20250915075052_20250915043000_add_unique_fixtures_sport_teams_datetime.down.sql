-- Drop composite unique constraint
ALTER TABLE fixtures
    DROP CONSTRAINT IF EXISTS fixtures_sport_teams_datetime_key;

