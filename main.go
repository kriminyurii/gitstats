package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

const (
	storePathName       string = "./.gitstatslocal" //TODO: Такое надо хранить где-то в .env не тут Переделать
	email               string = "yator0o+github@gmail.com"
	daysInLastSixMonths        = 183
)

func scanFolder(root string) error {
	var pathNames []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		regex, err := regexp.Compile(".git$")
		if err != nil {
			return err
		}
		if regex.MatchString(d.Name()) {
			pathNames = append(pathNames, path)
		}
		return nil
	})
	addGitPathToTheStore(storePathName, pathNames)
	return err
}

func addGitPathToTheStore(pathname string, newRepos []string) {
	repos := parseStoredFileLinesToSlice(pathname)
	repos = joinNewReposToSlice(newRepos, repos)
	dumbSliceToTheStoredFile(repos)
}

func joinNewReposToSlice(newRepos []string, existingRepos []string) []string {
	for _, v := range newRepos {
		if !contains(existingRepos, v) {
			existingRepos = append(existingRepos, v)
		}
	}
	return existingRepos
}

func dumbSliceToTheStoredFile(pathnames []string) {
	file, err := os.Create(storePathName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writter := bufio.NewWriter(file)
	for _, pathname := range pathnames {
		_, err := writter.WriteString(pathname + "\n")
		if err != nil {
			log.Fatalf("Error while writing pathname to the pathnames: %v", err)
		}
	}
	if err = writter.Flush(); err != nil {
		log.Fatal(err)
	}
}

func contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func parseStoredFileLinesToSlice(filepath string) []string {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filepath)
			if err != nil {
				panic(err)
			}
		} else {
			panic(fmt.Errorf("reading file error: %v", err))
		}
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		if err != io.EOF {
			panic(fmt.Errorf("reading file error: %v", err))
		}
	}
	return lines
}

func main() {
	err := scanFolder("../")
	if err != nil {
		log.Fatal(err)
	}
	proccessRepos(email)
	os.Exit(0)
}
