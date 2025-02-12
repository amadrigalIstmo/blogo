package cli

import (
	"blogo/internal/config"
	"errors"
	"fmt"
	"os"
)

// state contiene la configuración actual de la aplicación.
type state struct {
	cfg *config.Config
}

// command representa un comando ingresado por el usuario.
type command struct {
	name string
	args []string
}

// handlerLogin maneja el comando "login".
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("error: se requiere un nombre de usuario")
	}

	username := cmd.args[0]
	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("error al actualizar la configuración: %w", err)
	}

	fmt.Printf("Usuario actualizado correctamente a: %s\n", username)
	return nil
}

// commands almacena todos los comandos disponibles.
type commands struct {
	handlers map[string]func(*state, command) error
}

// register agrega un nuevo comando al sistema.
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// run ejecuta un comando si está registrado.
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("error: comando '%s' no reconocido", cmd.name)
	}
	return handler(s, cmd)
}

// RunCLI procesa los argumentos y ejecuta los comandos.
func RunCLI() {
	// Verificar argumentos mínimos
	if len(os.Args) < 2 {
		fmt.Println("error: no se proporcionaron suficientes argumentos")
		os.Exit(1)
	}

	// Leer la configuración
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error al leer la configuración: %v\n", err)
		os.Exit(1)
	}

	// Crear el estado de la aplicación
	appState := &state{cfg: &cfg}

	// Inicializar comandos
	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	// Extraer nombre del comando y argumentos
	cmd := command{name: os.Args[1], args: os.Args[2:]}

	// Ejecutar el comando
	if err := cmds.run(appState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
