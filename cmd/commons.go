package cmd

import (
	"bufio"
	"log"
	"os/user"
)

func getDotFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	pathToDotfile := usr.HomeDir + "/.gitgarden"
	return pathToDotfile
}

func extractExistingRepos(filePath string) []string {
	f := openOrCreate(filePath)
	defer f.Close()

	var lines []string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err) // No need to check for EOF, Scanner doesn't return io.EOF
	}

	return lines
}
