CREATE TABLE IF NOT EXISTS Team (
    team_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS Player (
    player_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INT,
    FOREIGN KEY (team_id) REFERENCES Team(team_id)
);

CREATE TABLE IF NOT EXISTS GameStats (
    game_id SERIAL PRIMARY KEY,
    player_id INT NOT NULL,
    points INT CHECK (points >= 0),
    rebounds INT CHECK (rebounds >= 0),
    assists INT CHECK (assists >= 0),
    steals INT CHECK (steals >= 0),
    blocks INT CHECK (blocks >= 0),
    fouls INT CHECK (fouls >= 0 AND fouls <= 6),
    turnovers INT CHECK (turnovers >= 0),
    minutes_played FLOAT CHECK (minutes_played >= 0.0 AND minutes_played <= 48.0),
    FOREIGN KEY (player_id) REFERENCES Player(player_id)
);