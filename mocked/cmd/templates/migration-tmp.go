package templates

const Migration = `-- +migrate Up
{{if .Create}}CREATE TABLE {{.TableName}} (
		id AUTO_INCREMENT PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
-- +migrate Down
DROP TABLE {{.TableName}};
	{{else if .Alter}}
		{{if .TableName}}
ALTER TABLE {{.TableName}}
-- +migrate Down
		{{end}}
	{{else}}
-- +migrate Up
//Write your SQL
-- +migrate Down
//Write your SQL
	{{end}}
`
