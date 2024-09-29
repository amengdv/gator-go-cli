package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/amengdv/gator/internal/config"
	"github.com/amengdv/gator/internal/database"
	_ "github.com/lib/pq"
)


func main() {

    conf, err := config.Read()
    if err != nil {
        log.Fatalf("Failed to read config: %v\n", err)
    }

    db, err := sql.Open("postgres", conf.DbUrl)
    if err != nil {
        log.Fatalf("Failed to connect to a database: %v\n", err)
    }

    dbQueries := database.New(db)

    s := state{
        db: dbQueries,
        cfg: &conf,
    }

    c := commands{
        handlers: make(map[string]func(*state, command) error),
    }

    c.register("login", handlerLogin)
    c.register("register", handlerRegister)
    c.register("reset", handlerReset)
    c.register("users", handlerUsers)
    c.register("agg", handlerAgg)
    c.register("addfeed", handlerAddFeed)
    c.register("feeds", handlerFeeds)
    c.register("follow", handlerFollow)
    c.register("following", handlerFollowing)

    argsWithoutProg := os.Args[1:]

    if len(argsWithoutProg) < 1 {
        log.Printf("Require a command name")
        os.Exit(1)
    }

    commandName := argsWithoutProg[0]
    commandArgs := argsWithoutProg[1:]

    err = c.run(&s, command{
        commandName,
        commandArgs,
    })

    if err != nil {
        log.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
