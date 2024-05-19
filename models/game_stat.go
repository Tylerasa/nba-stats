package models

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
