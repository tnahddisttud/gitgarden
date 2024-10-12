package cmd

import (
	"fmt"
	"sort"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const outOfRange int = 99999
const daysIn6Months int = 183
const weeksInLastSixMonths int = 26

type column []int

func Stats(user string) {
	commits := processRepos(user)
	printCommitGarden(commits)
}

func processRepos(user string) map[int]int {
	dotfile := getDotFilePath()
	repos := extractExistingRepos(dotfile)
	daysInMap := daysIn6Months

	commits := make(map[int]int, daysInMap)
	for i := daysInMap; i >= 0; i-- {
		commits[i] = 0
	}

	for _, path := range repos {
		commits = fillCommits(path, user, commits)
	}

	return commits
}

func fillCommits(path string, user string, commits map[int]int) map[int]int {
	repo, err := git.PlainOpen(path)
	if err != nil {
		panic(err)
	}

	ref, err := repo.Head()
	if err != nil {
		panic(err)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		panic(err)
	}

	offset := calculateOffset()

	err = iterator.ForEach(func(c *object.Commit) error {
		daysAgo := countDaysSince(c.Author.When) + offset
		if c.Author.Email != user {
			return nil
		}

		if daysAgo != outOfRange {
			commits[daysAgo]++
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return commits
}

func calculateOffset() int {
	var offset int
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		offset = 7
	case time.Monday:
		offset = 6
	case time.Tuesday:
		offset = 5
	case time.Wednesday:
		offset = 4
	case time.Thursday:
		offset = 3
	case time.Friday:
		offset = 2
	case time.Saturday:
		offset = 1
	}
	return offset
}

func getBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return startOfDay
}

func countDaysSince(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())
	for date.Before(now) {
		date = date.Add(time.Hour * 24)
		days++
		if days > daysIn6Months {
			return outOfRange
		}
	}
	return days
}

func printCommitGarden(commits map[int]int) {
	keys := sortMapsToSlices(commits)
	cols := buildColumns(keys, commits)
	printCells(cols)
}

func sortMapsToSlices(commits map[int]int) []int {
	var keys []int

	for k := range commits {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	return keys
}

func buildColumns(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	col := column{}

	for _, k := range keys {
		week := int(k / 7)
		dayInWeek := k % 7

		if dayInWeek == 0 {
			col = column{}
		}

		col = append(col, commits[k])

		if dayInWeek == 6 {
			cols[week] = col
		}
	}

	return cols
}

func printCells(cols map[int]column) {
	printMonths()

	for j := 6; j >= 0; j-- {
		for i := weeksInLastSixMonths + 1; i >= 0; i-- {
			if i == weeksInLastSixMonths+1 {
				printDayCols(j)
			}

			if col, ok := cols[i]; ok {
				if i == 0 && j == calculateOffset()-1 {
					printCell(col[j], true)
					continue
				} else {
					if len(col) > j {
						printCell(col[j], false)
						continue
					}
				}
			}

			printCell(0, false)
		}
		fmt.Printf("\n")
	}
}
