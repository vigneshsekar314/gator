package main

import (
	"github.com/vigneshsekar314/gator/internal/config"
	"log"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}
	state := GetNewState(&conf)
	cmds := GetNewCommands()
	Register(&cmds)
	if err := Run(&state, &cmds); err != nil {
		log.Fatal(err)
	}
}
