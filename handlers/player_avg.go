package handlers

import (
	"context"
	"encoding/json"
	"nbamodule/db"
	"nbamodule/models"
	"net/http"
)

func GetPlayerAverage(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Path[len("/average/player/"):]

	var avg models.GameStat

	rows, err := db.DB.Query(context.Background(), `
        SELECT
            player_id, AVG(points), AVG(rebounds), AVG(assists), AVG(steals), AVG(blocks), AVG(fouls), AVG(turnovers), AVG(minutes_played)
        FROM GameStats
        WHERE player_id = $1
        GROUP BY player_id
    `, playerID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&avg.PlayerID, &avg.Points, &avg.Rebounds, &avg.Assists, &avg.Steals, &avg.Blocks, &avg.Fouls, &avg.Turnovers, &avg.MinutesPlayed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(avg)
}
