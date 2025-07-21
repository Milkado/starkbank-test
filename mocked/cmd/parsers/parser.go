package parsers

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"test/starkbank/mocked/cmd/common"
)

func parseToFile(tmpfile string, folder string) *template.Template {
	//create template
	tmpl, err := template.New("migration").Parse(tmpfile)
	if err != nil {
		fmt.Println(common.Red, "Error parsing template: ")
		fmt.Println(err.Error(), common.Reset)
		return nil
	}

	//check if folder exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0600)
	}

	return tmpl
}

func parseController(data *controllerData, tmpfile string, folder string, filename string) {
	tmpl := parseToFile(tmpfile, folder)
	if tmpl == nil {
		return
	}

	var bytes bytes.Buffer
	err := tmpl.Execute(&bytes, data)
	if err != nil {
		fmt.Println(common.Red, "error executing template: ")
		fmt.Println(err.Error(), common.Reset)
		return
	}

	file := folder + "/" + filename
	// write from buffer
	err = os.WriteFile(file, bytes.Bytes(), 0600)
	if err != nil {
		fmt.Println(common.Red, "error writing file: ")
		fmt.Println(err.Error(), common.Reset)
		return
	}

	fmt.Println(common.Green, "Controller created successfully: ", filename, common.Reset)
}

func parseMigration(data *tmpData, tmpfile string, folder string, filename string) {
	tmpl := parseToFile(tmpfile, folder)
	if tmpl == nil {
		return
	}

	var bytes bytes.Buffer
	err := tmpl.Execute(&bytes, data)
	if err != nil {
		fmt.Println(common.Red, "error executing template: ")
		fmt.Println(err.Error(), common.Reset)
		return
	}

	file := folder + "/" + filename
	// write from buffer
	err = os.WriteFile(file, bytes.Bytes(), 0600)
	if err != nil {
		fmt.Println(common.Red, "error writing file: ")
		fmt.Println(err.Error(), common.Reset)
		return
	}

	fmt.Println(common.Green, "Migration created successfully: ", filename, common.Reset)
}
