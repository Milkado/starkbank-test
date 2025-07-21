package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"test/starkbank/mocked/cmd/common"
)

type (
	DB struct {
		db *sql.DB
	}

	MigrationFile struct {
		Name     string
		FullPath string
	}
)

func Executor(db *sql.DB) *DB {
	return &DB{db: db}
}

func (e *DB) Migrate() error {
	//Get absolute path
	absPath, err := filepath.Abs("./db/migrations")
	if err != nil {
		return err
	}

	//Scan for files
	files, err := scanMigrationFiles(absPath)
	if err != nil {
		return err
	}

	//Get pending migrations
	pending, err := getPendingMigrations(e.db, files)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		log.Println(common.Green, "Nothing to migrate", common.Reset)
		return nil
	}

	//Start transaction
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	//Current batch
	var batch int
	err = tx.QueryRow("select coalesce(max(batch), 0) from migrations").Scan(&batch)
	if err != nil {
		return err
	}
	batch++

	for _, file := range pending {
		//Execute migration
		log.Printf(common.Yellow+"Executing migration %s"+common.Reset, file.Name)

		//Extract up sql
		sqlContent, err := os.ReadFile(file.FullPath)
		if err != nil {
			return fmt.Errorf("error reading migration %s: %w", file.Name, err)
		}
		//Calculate checksum
		checksum := calculateCheckSum(sqlContent)

		up := getSqlContent(string(sqlContent), "Up")
		if up == "" {
			return fmt.Errorf(common.Red+"No up sql found in migration %s"+common.Reset, file.Name)
		}

		_, err = tx.Exec(string(up))
		if err != nil {
			return fmt.Errorf("error executing migration %s: %w", file.Name, err)
		}
		//Record migration
		_, err = tx.Exec("insert into migrations (batch, name, checksum) values (?, ?, ?)", batch, file.Name, checksum)
		if err != nil {
			return err
		}

		log.Printf(common.Green+"Migrated: %s"+common.Reset, file.Name)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf(common.Green+"Successfully migrated %d migrations\n"+common.Reset, len(pending))
	return nil
}

func (e *DB) Rollback() error {
	//Get last batch
	batch, err := getLastBatch(e.db)
	if err != nil {
		return err
	}
	if batch == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	//Get absolute path
	absPath, err := filepath.Abs("./db/migrations")
	if err != nil {
		return err
	}
	//Get last batch migrations
	migrations, err := getLastBatchMigrations(e.db, batch, absPath)
	if err != nil {
		return err
	}

	//Start transaction
	tx, err := e.db.Begin()
	if err != nil {
		return err
	}

	for _, file := range migrations {
		//Execute migration
		log.Printf(common.Yellow+"Rolling back migration %s"+common.Reset, file.Name)

		//Extract up sql
		sqlContent, err := os.ReadFile(file.FullPath)
		if err != nil {
			return fmt.Errorf("error reading migration %s: %w", file.Name, err)
		}

		down := getDown(string(sqlContent))
		if down == "" {
			return fmt.Errorf(common.Red+"No down sql found in migration %s"+common.Reset, file.Name)
		}

		_, err = tx.Exec(string(down))
		if err != nil {
			return fmt.Errorf("error rolling back migration %s: %w", file.Name, err)
		}
		//Delete migration record
		_, err = tx.Exec("delete from migrations where name = ?", file.Name)
		if err != nil {
			return err
		}

		log.Printf(common.Green+"Migration %s rolled back"+common.Reset, file.Name)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf(common.Green+"Successfully rolled back %d migrations\n"+common.Reset, len(migrations))
	return nil
}
