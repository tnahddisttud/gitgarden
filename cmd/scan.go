package cmd

import (
	"log"
	"os"
	"slices"
	"strings"
)

func Scan(path string) {
	print("Scanning this path: ", path)
	repositories := scanFolders(path)
	dotfilePath := getDotFilePath()
	saveReposInFile(repositories, dotfilePath)
}

func scanFolders(folder string) []string {
	return scanFoldersRecursively(make([]string, 30), folder)
}

func scanFoldersRecursively(folders []string, folder string) []string {
	ignored_folders := []string{".env", "node_modules", "vendor"}
	folder = strings.TrimSuffix(folder, "/")

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			path := folder + "/" + file.Name()

			if file.Name() == ".git" {
				path = strings.TrimSuffix(path, "/.git")
				folders = append(folders, path)
				continue
			}

			if slices.Contains(ignored_folders, file.Name()) {
				continue
			}

			folders = scanFoldersRecursively(folders, path)
		}
	}
	return folders
}


func saveReposInFile(repos []string, dotfile string) {
	existingRepos := extractExistingRepos(dotfile)
	allRepos := appendRepos(repos, existingRepos)
	updateReposInFile(allRepos, dotfile)
}

func openOrCreate(filePath string) *os.File {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	return f
}



func appendRepos(newRepos []string, existingRepos []string) []string {
	var result []string
	if existingRepos != nil {
		result = append([]string{}, existingRepos...)
	} else {
		result = make([]string, 0)
	}

	for _, repo := range newRepos {
		if !slices.Contains(result, repo) {
			result = append(result, repo)
		}
	}
	return result
}

func updateReposInFile(repos []string, dotfile string) {
	content := strings.Join(repos, "\n")
	os.WriteFile(dotfile, []byte(content), 0755)
}
