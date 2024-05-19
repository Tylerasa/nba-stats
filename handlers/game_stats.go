package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"nbamodule/db"
	"nbamodule/models"
	"net/http"
)

func CreateGameStat(w http.ResponseWriter, r *http.Request) {
	var stat models.GameStat
	err := json.NewDecoder(r.Body).Decode(&stat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exists bool
	err = db.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM Player WHERE player_id=$1)", stat.PlayerID).Scan(&exists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Player does not exist", http.StatusBadRequest)
		return
	}

	sqlStatement := `
    INSERT INTO GameStats (player_id, points, rebounds, assists, steals, blocks, fouls, turnovers, minutes_played)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING game_id`
	id := 0
	err = db.DB.QueryRow(context.Background(), sqlStatement, stat.PlayerID, stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Fouls, stat.Turnovers, stat.MinutesPlayed).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New game stat created with ID: %d", id)
}
