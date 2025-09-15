-- Add composite unique constraint to prevent duplicate fixtures
ALTER TABLE fixtures
    ADD CONSTRAINT fixtures_sport_teams_datetime_key
    UNIQUE (sport_id, team_id_1, team_id_2, date_time);

