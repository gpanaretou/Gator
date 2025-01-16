package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gpanaretou/Gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	available map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	switch name {
	case "login":
		c.available[name] = f
	}
}

func (c *commands) run(s *state, cmd command) error {
	_, ok := c.available[cmd.name]
	if !ok {
		err := fmt.Errorf("%v is not a command", cmd.name)
		return err
	}

	switch cmd.name {
	case "login":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("-- '%v' command not implemented", cmd.name)
		return err
	}

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		err := fmt.Errorf("login requires exactly 1 arguement")
		return err
	}

	user_name := cmd.args[0]
	err := s.cfg.SetUser(user_name)
	if err != nil {
		return err
	}

	fmt.Printf("> User is set to: %v", user_name)

	return nil
}

func main() {
	var s state
	cfg := config.Read()
	s.cfg = &cfg
	// user := "kokos"

	cmds := commands{
		available: make(map[string]func(*state, command) error),
	}
	cmds.available["login"] = nil
	cmds.register("login", handlerLogin)
	args := os.Args

	if len(args) < 2 {
		log.Fatal("was expecting at least one arguement")
	}

	var cmd command
	cmd.name = args[1]
	cmd.args = args[2:]

	err := cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
