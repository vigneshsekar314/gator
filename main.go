package main

import _ "github.com/lib/pq"

import (
	"database/sql"
	"github.com/vigneshsekar314/gator/internal/config"
	"github.com/vigneshsekar314/gator/internal/database"
	"log"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
		return
	}
	db, err := sql.Open("postgres", conf.DbUrl)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbQueries := database.New(db)
	state := GetNewState(&conf)
	state.db = dbQueries
	cmds := GetNewCommands()
	RegisterCommands(&cmds)
	if err := Run(&state, &cmds); err != nil {
		log.Fatal(err)
	}
}
