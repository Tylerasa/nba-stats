package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database!")
}

type GameStat struct {
	PlayerID      int     `json:"player_id"`
	Points        float64 `json:"points"`
	Rebounds      float64 `json:"rebounds"`
	Assists       float64 `json:"assists"`
	Steals        float64 `json:"steals"`
	Blocks        float64 `json:"blocks"`
	Fouls         float64 `json:"fouls"`
	Turnovers     float64 `json:"turnovers"`
	MinutesPlayed float64 `json:"minutes_played"`
}

type TeamStat struct {
	TeamID        int     `json:"team_id"`
	Points        float64 `json:"points"`
	Rebounds      float64 `json:"rebounds"`
	Assists       float64 `json:"assists"`
	Steals        float64 `json:"steals"`
	Blocks        float64 `json:"blocks"`
	Fouls         float64 `json:"fouls"`
	Turnovers     float64 `json:"turnovers"`
	MinutesPlayed float64 `json:"minutes_played"`
}

type Player struct {
	Name   string `json:"name"`
	TeamID int    `json:"team_id"`
}

type Team struct {
	Name string `json:"name"`
}

func createGameStat(w http.ResponseWriter, r *http.Request) {
	var stat GameStat
	err := json.NewDecoder(r.Body).Decode(&stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if player exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Player WHERE player_id=$1)", stat.PlayerID).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Player does not exist", http.StatusBadRequest)
		return
	}

	// Check if the game stat already exists
	// err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM GameStats WHERE player_id=$1 AND points=$2 AND rebounds=$3 AND assists=$4 AND steals=$5 AND blocks=$6 AND fouls=$7 AND turnovers=$8 AND minutes_played=$9)",
	// 	stat.PlayerID, stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Fouls, stat.Turnovers, stat.MinutesPlayed).Scan(&exists)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// if exists {
	// 	http.Error(w, "Game stat already exists", http.StatusConflict)
	// 	return
	// }

	sqlStatement := `
    INSERT INTO GameStats (player_id, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING game_id`
	id := 0
	err = db.QueryRow(sqlStatement, stat.PlayerID, stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Fouls, stat.Turnovers, stat.MinutesPlayed).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New game stat created with ID: %d", id)
}

func createPlayer(w http.ResponseWriter, r *http.Request) {
	var player Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the player already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Player WHERE name=$1 AND team_id=$2)", player.Name, player.TeamID).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Player already exists", http.StatusConflict)
		return
	}

	sqlStatement := `
    INSERT INTO Player (name, team_id)
    VALUES ($1, $2) RETURNING player_id`
	id := 0
	err = db.QueryRow(sqlStatement, player.Name, player.TeamID).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New player created with ID: %d", id)
}

func createTeam(w http.ResponseWriter, r *http.Request) {
	var team Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the team already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Team WHERE name=$1)", team.Name).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Team already exists", http.StatusConflict)
		return
	}

	sqlStatement := `
    INSERT INTO Team (name)
    VALUES ($1) RETURNING team_id`
	id := 0
	err = db.QueryRow(sqlStatement, team.Name).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New team created with ID: %d", id)
}

func main() {
	initDB()
	dbRun()
	defer db.Close()

	http.HandleFunc("/stats", createGameStat)
	http.HandleFunc("/players", createPlayer)
	http.HandleFunc("/teams", createTeam)
	http.HandleFunc("/average/player/", getPlayerAverage)
	http.HandleFunc("/average/team/", getTeamAverage)
	fmt.Println("Server started @ http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
