package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Load database URL
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	s := &state{
		config: &cfg,
		db:     database.New(db),
	}

	// Initialise commands map
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register commands inside
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerGetUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("feeds", handlerListAllFeeds)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))

	// Check if enough arguments were provided
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Not enough arguments. Please provide a command.")
		os.Exit(1)
	}

	// Create a command struct
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	// Run the command
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

}
