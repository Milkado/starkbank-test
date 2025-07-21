package db

import (
	"database/sql"
	"fmt"
	"log"
	"test/starkbank/mocked/cmd/common"

	"github.com/go-sql-driver/mysql"
)
var migrationTable = false
type DbConn struct {
	User   string
	Pass   string
	Addr   string
	DbName string
}

func Connect(conn DbConn) (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = conn.User
	cfg.Passwd = conn.Pass
	cfg.Net = "tcp"
	cfg.Addr = conn.Addr //127.0.0.1:3306
	cfg.DBName = conn.DbName
	cfg.ParseTime = true

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("error connecting to database", err.Error())
	}

	//Create migration table
	if !migrationTable {
		if exists, err := tableExists(conn, db); err != nil {
			return nil, err
		} else if !exists {
			err := createMigrationsTable(db)
			if err != nil {
				return nil, err
			}
		}
		migrationTable = true
	}

	return db, nil
}
func createMigrationsTable(db *sql.DB) error {
	// Add logging to debug
	fmt.Println(common.Yellow + "Creating migrations table..." + common.Reset)
	sql := `CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		batch INT NOT NULL,
		checksum TEXT NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("error creating migrations table, error: %v", err)
	}

	fmt.Println(common.Green + "Migrations table created successfully" + common.Reset)
	return nil
}

func tableExists(conn DbConn, db *sql.DB) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
        SELECT 1 FROM information_schema.tables 
        WHERE table_schema = ?
        AND table_name = 'migrations'
    )`
	err := db.QueryRow(query, conn.DbName).Scan(&exists)
	return exists, err
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

// query example
func albumsByArtist(name string, db *sql.DB) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}
