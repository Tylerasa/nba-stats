package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"nbamodule/db"
	"nbamodule/models"
	"net/http"
)

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var exists bool
	err = db.DB.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM Team WHERE name=$1)", team.Name).Scan(&exists)
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
	err = db.DB.QueryRow(context.Background(), sqlStatement, team.Name).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New team created with ID: %d", id)
}
