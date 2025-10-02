package config

import (
	"errors"
	"fmt"
	"os"
)

type state struct {
	config *Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	cmdList map[string]func(*state, command) error
}

func GetNewState(config *Config) state {
	return state{
		config: config,
	}
}

func GetNewCommands() commands {
	return commands{
		cmdList: make(map[string]func(*state, command) error),
	}
}

func Register(cmds *commands) {
	cmds.register("login", handlerLogin)
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
	if err := s.config.SetUser(cmd.arguments[0]); err != nil {
		return err
	}
	fmt.Printf("user has been set to %s\n", cmd.arguments[0])
	return nil
}
