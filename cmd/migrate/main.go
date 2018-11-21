package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/ykpythemind/funho/config"
)

func main() {
	config := config.Load()
	log.Println("migration " + config.APPName + "...")
	log.Println("mode " + config.Env)

	db, err := sql.Open("mysql", config.DBAddr())
	if err != nil {
		panic(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	var isDown bool
	if os.Getenv("DOWN") != "" {
		isDown = true
	}

	var n int
	if isDown {

		n, err = migrate.Exec(db, "mysql", migrations, migrate.Down)
	} else {
		n, err = migrate.Exec(db, "mysql", migrations, migrate.Up)

	}
	if err != nil {
		log.Println("migration Error!")
		panic(err)
	}
	fmt.Printf("Applied %d migrations!\n", n)
}
