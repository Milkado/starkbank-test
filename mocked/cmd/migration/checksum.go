package migration

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
)

func calculateCheckSum(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

func verifyIntegrity(db *sql.DB, files []MigrationFile) error {
	//Read file content
	for _, file := range files {
		

		content, err := os.ReadFile(file.FullPath)
		if err != nil {
			return fmt.Errorf("error reading migration %s: %w", file.Name, err)
		}

		checksum := calculateCheckSum(content)

		//Check againsta stored checksum
		var storedChecksum string
		err = db.QueryRow("select checksum from migrations where name = ?", file.Name).Scan(&storedChecksum)
		if err == sql.ErrNoRows {
			return nil
		}
		if err != nil {
			return err
		}
		if storedChecksum != checksum {
			return fmt.Errorf("checksum mismatch for migration %s", file.Name)
		}
	}
	return nil
}
