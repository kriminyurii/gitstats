package main

import (
	"fmt"
	"log"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func proccessRepos(email string) {
	const REFERENCE_NOT_FOUND_ERROR string = "reference not found"

	repos := parseStoredFileLinesToSlice(storePathName)
	fmt.Println(repos, "repos")
	commits := make(map[string]int)
	for _, repo := range repos {
		c, err := fillRepoInfo(repo, email, commits)
		if err != nil {
			if err.Error() == REFERENCE_NOT_FOUND_ERROR {
				fmt.Printf("There is no commits in this local repository: %s\n", repo)
				continue
			}
			log.Fatalf("fill repo error: %v", err)
		}
		commits = c
	}
	fmt.Println(commits, "commits")
}

func fillRepoInfo(repoPath, email string, commits map[string]int) (map[string]int, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	refHead, err := repo.Head()
	if err != nil {
		return commits, err
	}
	commitIterator, err := repo.Log(&git.LogOptions{From: refHead.Hash()})
	if err != nil {
		return nil, err
	}
	err = commitIterator.ForEach(func(c *object.Commit) error {
		if c.Author.Email != email {
			return nil
		} else {
			commits[email]++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commits, nil
}
