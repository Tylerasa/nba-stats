package main

import (
	"encoding/json"
	"net/http"
)

func getPlayerAverage(w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Path[len("/average/player/"):]

	var avg GameStat

	rows, err := db.Query(`
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

func getTeamAverage(w http.ResponseWriter, r *http.Request) {
	teamID := r.URL.Path[len("/average/team/"):]

	var avg TeamStat
	rows, err := db.Query(`
        SELECT Player.team_id, AVG(GameStats.points), AVG(GameStats.rebounds), AVG(GameStats.assists), AVG(GameStats.steals), AVG(GameStats.blocks), AVG(GameStats.fouls), AVG(GameStats.turnovers), AVG(GameStats.minutes_played)
        FROM GameStats
        JOIN Player ON GameStats.player_id = Player.player_id
        WHERE Player.team_id = $1
        GROUP BY Player.team_id
    `, teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&avg.TeamID, &avg.Points, &avg.Rebounds, &avg.Assists, &avg.Steals, &avg.Blocks, &avg.Fouls, &avg.Turnovers, &avg.MinutesPlayed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(avg)
}
