package main

import (
	"flag"

	lenslocked "lenslocked/controllers/appController"
)

func main() {
	configMessage := "Provide this flag in production." +
		" This ensures that a config.json file is provided " +
		"before the applications starts."
	configRequired := flag.Bool("prod", false, configMessage)
	flag.Parse()
	lenslocked.NewApp(*configRequired).Run()
}
