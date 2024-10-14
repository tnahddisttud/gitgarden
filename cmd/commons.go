package cmd

import (
	"log"
	"os/user"
	"bufio"
	"io"
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
		if err != io.EOF {
			panic(err)
		}
	}
	return lines
}
