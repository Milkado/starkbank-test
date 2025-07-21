package parsers

import (
	"strings"
	"test/starkbank/mocked/cmd/templates"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type controllerData struct {
	Name string
}

func GenerateControllerFile(name string) {
	//first letter of each word to uppercase
	name = cases.Title(language.English, cases.Compact).String(name)

	//remove spaces
	name = strings.ReplaceAll(name, " ", "")

	data := controllerData{
		Name: name,
	}

	filename := cases.Lower(language.English, cases.Compact).String(name) + ".go"
	parseController(&data, templates.Controller, "controllers", filename)
}
