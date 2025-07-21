package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"test/starkbank/mocked/cmd/common"
)

func scanMigrationFiles(directory string) ([]MigrationFile, error) {
	var migrations []MigrationFile

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		fileName := file.Name()
		trimmedname := strings.TrimSuffix(fileName, ".sql")

		filepath := filepath.Join(directory, fileName)

		migrations = append(migrations, MigrationFile{
			Name:     trimmedname,
			FullPath: filepath,
		})
	}

	return migrations, nil
}

func getPendingMigrations(db *sql.DB, files []MigrationFile) ([]MigrationFile, error) {
	var pending []MigrationFile

	//Check integrity before get pending
	err := verifyIntegrity(db, files)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		exists := false
		err := db.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM migrations WHERE name = ?)",
			file.Name).Scan(&exists)
		if err != nil {
			fmt.Println(common.Purple, err.Error(), common.Reset)
			return nil, err
		}

		if !exists {
			pending = append(pending, file)
		}
	}

	return pending, nil
}

func getSqlContent(content string, blocktype string) string {
	var sql string
	if blocktype == "Up" {
		sql = getUp(content)
	}
	if blocktype == "Down" {
		sql = getDown(content)
	}
	return sql
}

func getUp(content string) string {
	up := "-- +migrate Up"
	down := "-- +migrate Down"

	lines := strings.Split(content, "\n")
	var block []string
	var collecting bool
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == up {
			collecting = true
			continue
		}

		if trimmedLine == down {
			collecting = false
			continue
		}

		if collecting {
			block = append(block, line)
		}
	}

	return strings.TrimSpace(strings.Join(block, "\n"))
}

func getDown(content string) string {
	up := "-- +migrate Up"
	down := "-- +migrate Down"

	lines := strings.Split(content, "\n")
	var block []string
	var collecting bool
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == up {
			collecting = false
			continue
		}

		if trimmedLine == down {
			collecting = true
			continue
		}

		if collecting {
			block = append(block, line)
		}
	}

	return strings.TrimSpace(strings.Join(block, "\n"))
}

func getLastBatch(db *sql.DB) (int, error) {
	var batch int
	err := db.QueryRow("select coalesce(max(batch), 0) from migrations").Scan(&batch)
	if err != nil {
		return 0, err
	}
	return batch, nil
}

func getLastBatchMigrations(db *sql.DB, batch int, absPath string) ([]MigrationFile, error) {
	var migrations []MigrationFile

	rows, err := db.Query("select name, checksum from migrations where batch = ?", batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var checksum string
		err = rows.Scan(&name, &checksum)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, MigrationFile{
			Name:     name,
			FullPath: filepath.Join(absPath, name+".sql"),
		})
	}

	return migrations, nil
}
