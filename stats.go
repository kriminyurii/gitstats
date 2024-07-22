package main

import (
	"fmt"
	"log"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Offset struct {
	WeekDay int
	Row     int
	Commits int
}

func proccessRepos(email string) map[time.Time]Offset {
	const REFERENCE_NOT_FOUND_ERROR string = "reference not found"
	repos := parseStoredFileLinesToSlice(storePathName)
	commits := make(map[time.Time]Offset)
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
	return commits
}

func fillRepoInfo(repoPath, email string, commits map[time.Time]Offset) (map[time.Time]Offset, error) {
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
			offset := Offset{0, 0, 0}
			if prevOffset, ok := commits[date]; ok {
				offset.Commits = prevOffset.Commits + 1
			}
			offset.Row = getRowOffset(date)
			offset.WeekDay = getDayOffset(date)
			commits[date] = offset
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return commits, nil
}

func getDayOffset(date time.Time) int {
	return int(date.Weekday())
}

func getRowOffset(date time.Time) int {
	sixMonthsAgo := GetLastHalfYear().AddDate(0, 0, -time.Now().Day()+1)
	const daysInWeek int = 7
	const hoursInDay int = 24
	if date.Before(sixMonthsAgo) {
		return -1
	}
	durationSince := date.Sub(sixMonthsAgo)
	daysSinceSixMonthsAgo := int(durationSince.Hours() / float64(hoursInDay))
	return daysSinceSixMonthsAgo / daysInWeek
}

func GetLastHalfYear() time.Time {
	return time.Now().AddDate(0, 0, -daysInLastSixMonths)
}

func GetLastHalfYearInMonths() []string {
	sixMonthsAgo := GetLastHalfYear().AddDate(0, 0, -time.Now().Day()+1)

	var lastSixMonths []string
	for month := sixMonthsAgo; month.Before(time.Now()); month = month.AddDate(0, 1, 0) {
		lastSixMonths = append(lastSixMonths, month.Format("Jan"))
	}
	return lastSixMonths
}

func filterCommitByDate(commits map[time.Time]Offset) map[time.Time]Offset {
	filteredCommits := make(map[time.Time]Offset)
	halfYearAgo := GetLastHalfYear()

	for commitDate, offset := range commits {
		if commitDate.After(halfYearAgo) {
			filteredCommits[commitDate] = offset
		}
	}

	return filteredCommits
}
