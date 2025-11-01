package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jkk290/gator/internal/config"
	"github.com/jkk290/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	appState := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		addedCmds: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)

	userInput := os.Args
	if len(userInput) < 2 {
		fmt.Println("error not enough arguments passed")
		os.Exit(1)
	}

	cmdName := userInput[1]
	cmdArgs := userInput[2:]
	commandInstance := command{
		Name: cmdName,
		Args: cmdArgs,
	}

	if err := cmds.run(&appState, commandInstance); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
