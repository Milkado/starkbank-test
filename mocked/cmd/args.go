package main

import (
	"fmt"
	"os"
	"test/starkbank/helpers"
	"test/starkbank/mocked/cmd/common"
	"test/starkbank/mocked/cmd/migration"
	"test/starkbank/mocked/cmd/parsers"
	"test/starkbank/mocked/db"
)


var conn = db.DbConn{
	User: helpers.Env("DB_USER"),
	Pass: helpers.Env("DB_PASSWORD"),
	Addr: helpers.Env("DB_HOST") + ":" + helpers.Env("DB_PORT"),
	DbName: helpers.Env("DB_NAME"),
}

func command(cmd string) {
	switch cmd {
	case "create:migration":
		migrationCmd()
	case "create:controller":
		controllerCmd()
	case "migrate":
		migrateCmd()
	case "migrate:rollback":
		rollbackCmd()
	default:
		errorC()
	}
}

func migrationCmd() {
	if len(os.Args) < 3 {
		nameMissing()
		return
	}

	name := os.Args[2]

	parsers.GenerateMigrationFile(name)
}

func controllerCmd() {
	if len(os.Args) < 3 {
		nameMissing()
		return
	}

	name := os.Args[2]

	parsers.GenerateControllerFile(name)
}

func nameMissing() {
	fmt.Println(common.Red, "Name not specified.", common.Reset)
	fmt.Println(common.Yellow, " Use ./gomd command <name>", common.Reset)
}

func migrateCmd() {
	db, err := db.Connect(conn)
	if err != nil {
		fmt.Println(common.Red, "Error connecting to database:", err.Error(), common.Reset)
		return
	}
	err = migration.Executor(db).Migrate()
	if err != nil {
		fmt.Println(common.Red, "Error migrating:", err.Error(), common.Reset)
		return
	}
}

func rollbackCmd() {
	db, err := db.Connect(conn)
	if err != nil {
		fmt.Println(common.Red, "Error connecting to database:", err.Error(), common.Reset)
		return
	}
	err = migration.Executor(db).Rollback()
	if err != nil {
		fmt.Println(common.Red, "Error rolling back:", err.Error(), common.Reset)
		return
	}
}
