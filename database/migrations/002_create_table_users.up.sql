CREATE TABLE users(
    user_id VARCHAR(255) PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL,
    team_name uuid NOT NULL REFERENCES teams(team_id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL
);