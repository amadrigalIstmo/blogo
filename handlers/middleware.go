package handlers

import (
	"blogo/internal/database"
	"context"
	"errors"
	"fmt"
	"log"
)

// LoggedInMiddleware verifica si el usuario está autenticado antes de ejecutar el handler.
func LoggedInMiddleware(app Application, handler func(user database.User, args []string) error) func(args []string) error {
	return func(args []string) error {
		user, err := getCurrentUserMiddleWare(app)
		if err != nil {
			return errors.New("error: debes estar autenticado para ejecutar este comando")
		}
		return handler(user, args)
	}
}

// getCurrentUser obtiene el usuario autenticado desde la configuración y la base de datos.
func getCurrentUserMiddleWare(app Application) (database.User, error) {
	currentUser := app.GetConfig().CurrentUser
	if currentUser == "" {
		return database.User{}, fmt.Errorf("no hay usuario autenticado")
	}

	user, err := app.GetDB().GetUser(context.Background(), currentUser)
	if err != nil {
		log.Println("Error obteniendo usuario autenticado:", err)
		return database.User{}, fmt.Errorf("usuario no encontrado en la base de datos")
	}
	return user, nil
}
