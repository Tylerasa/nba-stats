package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitDB() {

	dotErr := godotenv.Load()
	if dotErr != nil {
		fmt.Println("Error loading .env file")
	}

	createTables := `
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
    );`

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	DB, err = pgxpool.Connect(context.Background(), connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, execErr := DB.Exec(context.Background(), createTables)
	if execErr != nil {
		log.Fatal(execErr)
	}

	fmt.Println("Tables created successfully!")

	fmt.Println("Successfully connected to the database!")
}
