ALTER TABLE teams ADD COLUMN api_id INT;

CREATE INDEX idx_teams_api_id ON teams (api_id);