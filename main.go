package main

import (
	"database/sql"

	"github.com/DL1793/gator/internal/database"
	_ "github.com/lib/pq"
)

import (
	"fmt"
	"log"
	"os"

	"github.com/DL1793/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	cmds := commands{
		make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerGetFeedFollowsForUser))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	args := os.Args

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	st := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	if len(args) < 2 {
		fmt.Println("no command specified")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := args[2:]

	command := command{
		cmdName,
		cmdArgs[:],
	}

	err = cmds.run(&st, command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
