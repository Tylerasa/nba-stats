package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"nbamodule/db"
	"nbamodule/models"
	"net/http"
)

func CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exists bool
	err = db.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM Player WHERE name=$1 AND team_id=$2)", player.Name, player.TeamID).Scan(&exists)
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
	err = db.DB.QueryRow(context.Background(), sqlStatement, player.Name, player.TeamID).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New player created with ID: %d", id)
}
