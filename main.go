package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
)

func scanFolder(root string) {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		regex, err := regexp.Compile(".git$")
		if err != nil {
			panic(err)
		}
		if regex.MatchString(d.Name()) {
			fmt.Println(d.Name(), "name")
		}
		return nil
	})
	log.Fatal(err)
}

func main() {
	scanFolder("../docker-nodejs-sample")
}
