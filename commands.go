package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/vigneshsekar314/gator/internal/config"
	"github.com/vigneshsekar314/gator/internal/database"
	"os"
	"time"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmdList map[string]func(*state, command) error
}

func GetNewState(config *config.Config) state {
	return state{
		config: config,
	}
}

func GetNewCommands() commands {
	return commands{
		cmdList: make(map[string]func(*state, command) error),
	}
}

func RegisterCommands(cmds *commands) {
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
}

func Run(state *state, cmds *commands) error {
	args := os.Args
	if len(args) < 2 {
		return errors.New("Not enough arguments")
	}
	cmd := command{}
	cmd.name = args[1]
	cmd.arguments = args[2:]
	if err := cmds.run(state, cmd); err != nil {
		return err
	}
	return nil
}

func (cmds *commands) run(s *state, cmd command) error {
	if s == nil || s.config == nil {
		return errors.New("state does not exists")
	}
	fnToRun, ok := cmds.cmdList[cmd.name]
	if !ok {
		return errors.New("command is not registered")
	}
	if err := fnToRun(s, cmd); err != nil {
		return err
	}
	return nil
}

func (cmds *commands) register(name string, fn func(*state, command) error) {
	cmds.cmdList[name] = fn
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Login command expects username")
	}
	usrNm := cmd.arguments[0]
	usr, err := s.db.GetUser(context.Background(), usrNm)
	if err != nil {
		return err
	}
	if usr.Name != usrNm {
		return errors.New("User is not found in the database")
	}
	if err := s.config.SetUser(usrNm); err != nil {
		return err
	}
	fmt.Printf("user has been set to %s\n", usrNm)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("Register command exxpects username")
	}
	timeNow := time.Now()
	usrNm := cmd.arguments[0]
	usr, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Name:      usrNm,
	})
	if err != nil {
		return err
	}
	if err := s.config.SetUser(usrNm); err != nil {
		return err
	}
	fmt.Printf("User is created. User ID: %v, createdAt: %v, updatedAt: %v, Name: %s\n", usr.ID, usr.CreatedAt, usr.UpdatedAt, usr.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUser(context.Background()); err != nil {
		return err
	}
	fmt.Println("All users have been deleted")
	return nil
}
