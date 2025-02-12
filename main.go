package main

import (
	"database/sql"
	"log"
	"os"

	"blogo/handlers"
	"blogo/internal/config"
	"blogo/internal/database"

	_ "github.com/lib/pq"
)

type application struct {
	db  *database.Queries
	cfg *config.Config
}

func (app *application) GetDB() *database.Queries {
	return app.db
}

func (app *application) GetConfig() *config.Config {
	return app.cfg
}

func (app *application) SaveConfig() error {
	return app.cfg.Save()
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	app := &application{
		db:  database.New(db),
		cfg: &cfg,
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: go run . [register|login|reset|users|agg] [arguments]")
	}

	switch os.Args[1] {
	case "reset":
		handlers.ResetDatabase(app)
	case "register":
		handlers.RegisterUser(app, os.Args[2:])
	case "login":
		handlers.LoginUser(app, os.Args[2:])
	case "users":
		handlers.ListUsers(app)
	case "agg":
		handlers.HandleAgg(app)
	default:
		log.Fatal("Invalid command")
	}
}
