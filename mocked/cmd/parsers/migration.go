package parsers

import (
	"fmt"
	"regexp"
	"strings"
	"test/starkbank/mocked/cmd/common"
	"test/starkbank/mocked/cmd/templates"
	"time"
)

type tmpData struct {
	TableName  *string
	Timestramp string
	Create     bool
	Alter      bool
}

func validateMigrationName(name string) error {
	re := `^[a-z0-9]+(_[a-z0-9]+)*$`
	matched, err := regexp.MatchString(re, name)
	if err != nil {
		return err
	}

	if !matched {
		return fmt.Errorf("invalid migration name. Use only lowercase letters, numbers and underscores or space as separators. Example: create_users_table or create users table")
	}

	return nil
}

func GenerateMigrationFile(name string) {
	name = strings.ReplaceAll(name, " ", "_")
	err := validateMigrationName(name)
	if err != nil {
		fmt.Println(common.Red, "Error:")
		fmt.Println(err.Error(), common.Reset)
		return
	}

	timestamp := time.Now().Format("20060102150405")

	parsedname := strings.ReplaceAll(name, "_", "")
	var tablename string
	var create, alter bool
	if strings.HasPrefix(name, "create") {
		tablename, create, alter = createData(parsedname)
	}

	data := tmpData{
		TableName:  &tablename,
		Timestramp: timestamp,
		Create:     create,
		Alter:      alter,
	}

	//create migration
	filename := timestamp + "_" + name + ".sql"

	parseMigration(&data, templates.Migration, "db/migrations", filename)

}

func createData(name string) (string, bool, bool) {
	create := true
	alter := false
	tablename := strings.ReplaceAll(name, "create", "")
	if strings.HasSuffix(tablename, "table") {
		tablename = strings.ReplaceAll(tablename, "table", "")
	}

	//remove first ans last underscores
	tablename = strings.Trim(tablename, "_")

	return tablename, create, alter
}
