package cmd

import (
	"fmt"
	"log"
	"sort"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	outOfRange            = 99999
	daysInSixMonths       = 183
	weeksInLastSixMonths  = 26
)

type column []int

func Stats(user string) {
	commits := processRepos(user)
	printCommitGarden(commits)
}

func processRepos(user string) map[int]int {
	repos := extractExistingRepos(getDotFilePath())
	commits := make(map[int]int, daysInSixMonths)

	for i := 0; i <= daysInSixMonths; i++ {
		commits[i] = 0
	}

	for _, path := range repos {
		fillCommits(path, user, commits)
	}

	return commits
}

func fillCommits(path, user string, commits map[int]int) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	ref, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}

	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}

	offset := calculateOffset()

	iterator.ForEach(func(c *object.Commit) error {
		if c.Author.Email != user {
			return nil
		}
		daysAgo := countDaysSince(c.Author.When) + offset
		if daysAgo != outOfRange {
			commits[daysAgo]++
		}
		return nil
	})
}

func calculateOffset() int {
	switch time.Now().Weekday() {
	case time.Sunday:
		return 7
	case time.Monday:
		return 6
	case time.Tuesday:
		return 5
	case time.Wednesday:
		return 4
	case time.Thursday:
		return 3
	case time.Friday:
		return 2
	default:
		return 1
	}
}

func getBeginningOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func countDaysSince(date time.Time) int {
	days := 0
	now := getBeginningOfDay(time.Now())

	for date.Before(now) {
		date = date.Add(24 * time.Hour)
		days++
		if days > daysInSixMonths {
			return outOfRange
		}
	}

	return days
}

func printCommitGarden(commits map[int]int) {
	keys := sortedKeys(commits)
	cols := buildColumns(keys, commits)
	printCells(cols)
}

func sortedKeys(commits map[int]int) []int {
	keys := make([]int, 0, len(commits))
	for k := range commits {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

func buildColumns(keys []int, commits map[int]int) map[int]column {
	cols := make(map[int]column)
	var col column

	for _, k := range keys {
		week := k / 7
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
				printCell(col, j, i == 0 && j == calculateOffset()-1)
			} else {
				printCell(nil, j, false)
			}
		}
		fmt.Println()
	}
}

func printMonths() {
	week := getBeginningOfDay(time.Now()).Add(-time.Duration(daysInSixMonths) * 24 * time.Hour)
	month := week.Month()
	fmt.Print("         ")

	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Print("    ")
		}

		week = week.Add(7 * 24 * time.Hour)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Println()
}

func printDayCols(day int) {
	switch day {
	case 1:
		fmt.Print(" Mon ")
	case 3:
		fmt.Print(" Wed ")
	case 5:
		fmt.Print(" Fri ")
	default:
		fmt.Print("     ")
	}
}

func printCell(col column, day int, today bool) {
	var escape string
	if col == nil || len(col) <= day {
		escape = "\033[48;5;235;38;5;235m  - "
	} else {
		val := col[day]
		switch {
		case today:
			escape = "\033[48;5;167;38;5;229m"
		case val == 0:
			escape = "\033[48;5;235;38;5;235m"
		case val <= 3:
			escape = "\033[48;5;142;38;5;235m"
		case val <= 6:
			escape = "\033[48;5;108;38;5;235m"
		case val <= 9:
			escape = "\033[48;5;109;38;5;235m"
		default:
			escape = "\033[48;5;175;38;5;235m"
		}
		fmt.Printf(escape+" %2d \033[0m", val)
		return
	}
	fmt.Print(escape + "\033[0m")
}

