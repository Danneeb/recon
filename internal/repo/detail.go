package repo

import (
	"sort"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Commit struct {
	Hash        string
	Message     string
	Author      string
	AuthorEmail string
	When        time.Time
	Date        string
}

type BarEntry struct {
	Label string
	Count int
	Pct   int // 0–100, relative to max in set
}

type RepoDetail struct {
	*Repo
	Commits         []Commit
	Contributors    []string
	LastCommitDate  string
	ReadMePath      string
	CommitsByMonth  []BarEntry
	CommitsByAuthor []BarEntry
	CommitsByDay    []BarEntry
	CommitCount     int
}

func GetRepoDetail(repo *Repo, path string) (*RepoDetail, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	iter, err := r.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	var commits []Commit
	seen := make(map[string]struct{})
	monthCounts := make(map[string]int)
	authorCounts := make(map[string]int)
	var dayCounts [7]int

	err = iter.ForEach(func(c *object.Commit) error {
		commits = append(commits, Commit{
			Hash:        c.Hash.String(),
			Message:     c.Message,
			Author:      c.Author.Name,
			AuthorEmail: c.Author.Email,
			When:        c.Author.When,
			Date:        c.Author.When.Format("2 Jan 2006"),
		})
		seen[c.Author.Name] = struct{}{}
		monthCounts[c.Author.When.Format("2006-01")]++
		authorCounts[c.Author.Name]++
		dayCounts[c.Author.When.Weekday()]++
		return nil
	})
	if err != nil {
		return nil, err
	}

	contributors := make([]string, 0, len(seen))
	for name := range seen {
		contributors = append(contributors, name)
	}

	// Commits by month — full history
	now := time.Now()
	numMonths := 1
	if len(commits) > 0 {
		oldest := commits[len(commits)-1].When
		numMonths = (now.Year()-oldest.Year())*12 + int(now.Month()) - int(oldest.Month()) + 1
	}
	months := make([]BarEntry, numMonths)
	maxMonth := 1
	for i := 0; i < numMonths; i++ {
		t := now.AddDate(0, -(numMonths-1-i), 0)
		count := monthCounts[t.Format("2006-01")]
		months[i] = BarEntry{Label: t.Format("Jan '06"), Count: count}
		if count > maxMonth {
			maxMonth = count
		}
	}
	for i := range months {
		months[i].Pct = months[i].Count * 100 / maxMonth
	}

	// Commits by author — top 10
	type kv struct {
		k string
		v int
	}
	authorsSorted := make([]kv, 0, len(authorCounts))
	for k, v := range authorCounts {
		authorsSorted = append(authorsSorted, kv{k, v})
	}
	sort.Slice(authorsSorted, func(i, j int) bool {
		return authorsSorted[i].v > authorsSorted[j].v
	})
	maxAuthor := 1
	if len(authorsSorted) > 0 {
		maxAuthor = authorsSorted[0].v
	}
	authorEntries := make([]BarEntry, len(authorsSorted))
	for i := range authorsSorted {
		authorEntries[i] = BarEntry{
			Label: authorsSorted[i].k,
			Count: authorsSorted[i].v,
			Pct:   authorsSorted[i].v * 100 / maxAuthor,
		}
	}

	// Commits by day of week (Mon–Sun)
	dayOrder := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
	dayLabels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	maxDay := 1
	for _, d := range dayCounts {
		if d > maxDay {
			maxDay = d
		}
	}
	dayEntries := make([]BarEntry, 7)
	for i, wd := range dayOrder {
		dayEntries[i] = BarEntry{
			Label: dayLabels[i],
			Count: dayCounts[wd],
			Pct:   dayCounts[wd] * 100 / maxDay,
		}
	}

	lastCommitDate := ""
	if len(commits) > 0 {
		lastCommitDate = commits[0].Date
	}

	return &RepoDetail{
		Repo:            repo,
		Commits:         commits,
		Contributors:    contributors,
		CommitCount:     len(commits),
		LastCommitDate:  lastCommitDate,
		CommitsByMonth:  months,
		CommitsByAuthor: authorEntries,
		CommitsByDay:    dayEntries,
	}, nil
}
