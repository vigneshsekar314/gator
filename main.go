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
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddfeed)
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", handlerFollow)
	cmds.register("following", handlerFollowing)
	if err := Run(&state, &cmds); err != nil {
		log.Fatal(err)
	}
}
