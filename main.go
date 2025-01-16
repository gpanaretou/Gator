package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gpanaretou/Gator/internal/config"
	"github.com/gpanaretou/Gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
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
	case "register":
		c.available[name] = f
	case "reset":
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
	case "register":
		err := c.available[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	case "reset":
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
		err := fmt.Errorf("login requires exactly 1 arguement: login <name>")
		return err
	}

	user_name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), user_name)
	if err != nil {
		return fmt.Errorf("%v user does not exist", user_name)
	}

	err = s.cfg.SetUser(user_name)
	if err != nil {
		return err
	}

	fmt.Printf("> User is set to: %v\n", user_name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("register requires 1 arguement: register <name>")
	}

	user_name := cmd.args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      user_name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return fmt.Errorf("user with name %v already exists", user_name)
	}

	s.cfg.SetUser(user_name)
	fmt.Printf("*SYSTEM: USER %v was created successfully\n", user_name)
	fmt.Printf("*SYSTEM: USER: \n%v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("reset requires no arguements")
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("*SYSTEM: SUCCESSFULLY RESET DATABASE")
	return nil
}

func main() {
	var s state
	cfg := config.Read()
	s.cfg = &cfg

	cmds := commands{
		available: make(map[string]func(*state, command) error),
	}
	cmds.available["login"] = nil
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	args := os.Args

	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		log.Fatal("error while trying to connect to the Database", err)
	}

	dbQueries := database.New(db)
	s.db = dbQueries

	if len(args) < 2 {
		log.Fatal("was expecting at least one arguement")
	}

	var cmd command
	cmd.name = args[1]
	cmd.args = args[2:]

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
