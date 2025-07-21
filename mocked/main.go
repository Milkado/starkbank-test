package main

import (
	"test/starkbank/mocked/routes"
)



func main() {
	e := routes.Api()

	e.Logger.Fatal(e.Start(":9090"))
}
