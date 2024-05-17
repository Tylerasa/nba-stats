package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func dbRun() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	dbHost, dbPort, dbUser, dbPassword, dbName)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	createTables(db)

	fmt.Println("Tables checked/created successfully!")
}

func createTables(db *sql.DB) {
	createTeamTable := `
    CREATE TABLE IF NOT EXISTS Team (
        team_id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL
    );`

	createPlayerTable := `
    CREATE TABLE IF NOT EXISTS Player (
        player_id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        team_id INT,
        FOREIGN KEY (team_id) REFERENCES Team(team_id)
    );`

	createGameStatsTable := `
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

	_, err := db.Exec(createTeamTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPlayerTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createGameStatsTable)
	if err != nil {
		log.Fatal(err)
	}
}
