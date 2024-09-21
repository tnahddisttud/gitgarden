package main

import (
	"flag"

	"github.com/tnahddisttud/gitgarden/cmd"
)

func main() {
	var path string
	var email string

	flag.StringVar(&path, "add", "", "add a folder to be tracked by gitgarden for any git repos")
	flag.StringVar(&email, "email", "", "email of the user to be scanned")
	flag.Parse()

	if path != "" {
		cmd.Scan(path)
		return
	}

	if email == "" {
		panic("`-email` cannot be empty!")
	}
	cmd.Stats(email)
}
