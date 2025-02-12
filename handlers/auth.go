package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"blogo/internal/config"
	"blogo/internal/database"

	"github.com/google/uuid"
)

type Application interface {
	GetDB() *database.Queries
	GetConfig() *config.Config
	SaveConfig() error
}

func RegisterUser(app Application, args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: register [username]")
	}

	username := args[0]

	_, err := app.GetDB().CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      username,
	})

	if err != nil {
		handleDBError(err, username)
	}

	app.GetConfig().CurrentUser = username
	if err := app.SaveConfig(); err != nil {
		log.Fatal("Error saving config:", err)
	}
	fmt.Printf("User %s registered\n", username)
}

func LoginUser(app Application, args []string) {
	if len(args) < 1 {
		log.Fatal("Usage: login [username]")
	}

	username := args[0]

	_, err := app.GetDB().GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("User %s does not exist", username)
	}

	app.GetConfig().CurrentUser = username
	if err := app.SaveConfig(); err != nil {
		log.Fatal("Error saving config:", err)
	}
	fmt.Printf("Welcome %s\n", username)
}

// Handler para resetear la base de datos
func ResetDatabase(app Application) {
	err := app.GetDB().ResetDatabase(context.Background()) // <-- Solo 1 valor de retorno
	if err != nil {
		log.Fatal("Error resetting database:", err)
	}
	fmt.Println("Database reset successfully")
}

func handleDBError(err error, username string) {
	if strings.Contains(err.Error(), "unique constraint") {
		log.Fatalf("User '%s' already exists", username)
	}
	log.Fatalf("Database error: %v", err)
}
