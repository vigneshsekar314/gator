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
	state := config.GetNewState(&conf)
	cmds := config.GetNewCommands()
	config.Register(&cmds)
	if err := config.Run(&state, &cmds); err != nil {
		log.Fatal(err)
	}

	// conf, err = config.Read()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("DB URL: %s\n", conf.DbUrl)
	// fmt.Printf("Username: %s\n", conf.CurrentUserName)
}
