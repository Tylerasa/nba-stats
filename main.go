package main

import (
	"context"
	"fmt"
	"log"
	"nbamodule/db"
	"nbamodule/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	db.InitDB()

	defer db.DB.Close()

	http.HandleFunc("/stats", handlers.CreateGameStat)
	http.HandleFunc("/players", handlers.CreatePlayer)
	http.HandleFunc("/teams", handlers.CreateTeam)
	http.HandleFunc("/average/player/", handlers.GetPlayerAverage)
	http.HandleFunc("/average/team/", handlers.GetTeamAverage)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	fmt.Println("Server started @ http://localhost:8080")

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
