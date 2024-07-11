package main

import (
	"fmt"
	"log"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func proccessRepos(email string) {
	const REFERENCE_NOT_FOUND_ERROR string = "reference not found"

	repos := parseStoredFileLinesToSlice(storePathName)
	commits := make(map[time.Time]int)
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
	commits = filterCommitByDate(commits)
	fmt.Println(commits, "commits")
}

func fillRepoInfo(repoPath, email string, commits map[time.Time]int) (map[time.Time]int, error) {
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
			date := c.Author.When.Truncate(24 * time.Hour)
			commits[date]++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func filterCommitByDate(commits map[time.Time]int) map[time.Time]int {
	filteredCommits := make(map[time.Time]int)
	halfYearAgo := time.Now().AddDate(0, 0, -daysInLastSixMonths)

	for commitDate, count := range commits {
		if commitDate.After(halfYearAgo) {
			filteredCommits[commitDate] = count
		}
	}

	return filteredCommits
}
