package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type commitSpec struct {
	when    time.Time
	author  string
	email   string
	message string
}

func makeRepoWithCommits(t *testing.T, specs []commitSpec) *Repo {
	t.Helper()
	dir := t.TempDir()

	r, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("PlainInit: %v", err)
	}
	w, err := r.Worktree()
	if err != nil {
		t.Fatalf("Worktree: %v", err)
	}

	filePath := filepath.Join(dir, "file.txt")
	for i, spec := range specs {
		if err = os.WriteFile(filePath, []byte(fmt.Sprintf("content %d", i)), 0644); err != nil {
			t.Fatal(err)
		}
		if _, err = w.Add("file.txt"); err != nil {
			t.Fatalf("Add: %v", err)
		}
		_, err = w.Commit(spec.message, &git.CommitOptions{
			Author: &object.Signature{
				Name:  spec.author,
				Email: spec.email,
				When:  spec.when,
			},
		})
		if err != nil {
			t.Fatalf("Commit: %v", err)
		}
	}

	return &Repo{Name: filepath.Base(dir), Path: dir, Branch: "master"}
}

func TestGetRepoDetailInvalidPath(t *testing.T) {
	r := &Repo{Name: "missing", Path: "/nonexistent/path", Branch: "main"}
	_, err := GetRepoDetail(r, r.Path)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestGetRepoDetailCommitCount(t *testing.T) {
	now := time.Now()
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "first", when: now.AddDate(0, -2, 0)},
		{author: "Alice", email: "alice@example.com", message: "second", when: now.AddDate(0, -1, 0)},
		{author: "Alice", email: "alice@example.com", message: "third", when: now},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	if detail.CommitCount != 3 {
		t.Errorf("expected CommitCount 3, got %d", detail.CommitCount)
	}
	if len(detail.Commits) != 3 {
		t.Errorf("expected 3 commits in slice, got %d", len(detail.Commits))
	}
}

func TestGetRepoDetailContributors(t *testing.T) {
	now := time.Now()
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "a1", when: now.AddDate(0, -2, 0)},
		{author: "Bob", email: "bob@example.com", message: "b1", when: now.AddDate(0, -1, 0)},
		{author: "Alice", email: "alice@example.com", message: "a2", when: now},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	if len(detail.Contributors) != 2 {
		t.Errorf("expected 2 unique contributors, got %d", len(detail.Contributors))
	}
}

func TestGetRepoDetailLastCommitDate(t *testing.T) {
	when := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "commit", when: when},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	want := "15 Jun 2024"
	if detail.LastCommitDate != want {
		t.Errorf("expected LastCommitDate %q, got %q", want, detail.LastCommitDate)
	}
}

func TestGetRepoDetailCommitDateFormat(t *testing.T) {
	when := time.Date(2024, 3, 5, 0, 0, 0, 0, time.UTC)
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "commit", when: when},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	want := "5 Mar 2024"
	if detail.Commits[0].Date != want {
		t.Errorf("expected Commit.Date %q, got %q", want, detail.Commits[0].Date)
	}
}

func TestGetRepoDetailCommitsByAuthor(t *testing.T) {
	now := time.Now()
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "a1", when: now.AddDate(0, 0, -3)},
		{author: "Alice", email: "alice@example.com", message: "a2", when: now.AddDate(0, 0, -2)},
		{author: "Alice", email: "alice@example.com", message: "a3", when: now.AddDate(0, 0, -1)},
		{author: "Bob", email: "bob@example.com", message: "b1", when: now},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	if len(detail.CommitsByAuthor) == 0 {
		t.Fatal("CommitsByAuthor is empty")
	}
	top := detail.CommitsByAuthor[0]
	if top.Label != "Alice" {
		t.Errorf("expected top author Alice, got %q", top.Label)
	}
	if top.Count != 3 {
		t.Errorf("expected Alice count 3, got %d", top.Count)
	}
	if top.Pct != 100 {
		t.Errorf("expected top author Pct 100, got %d", top.Pct)
	}
}

func TestGetRepoDetailCommitsByMonth(t *testing.T) {
	twoMonthsAgo := time.Now().AddDate(0, -2, 0)
	thisMonth := time.Now()
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "old", when: twoMonthsAgo},
		{author: "Alice", email: "alice@example.com", message: "new", when: thisMonth},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	if len(detail.CommitsByMonth) < 3 {
		t.Errorf("expected at least 3 month entries for commits 2 months apart, got %d", len(detail.CommitsByMonth))
	}
	// First entry is the oldest commit's month
	if detail.CommitsByMonth[0].Count != 1 {
		t.Errorf("expected first month count 1, got %d", detail.CommitsByMonth[0].Count)
	}
	// Last entry is this month
	last := detail.CommitsByMonth[len(detail.CommitsByMonth)-1]
	if last.Count != 1 {
		t.Errorf("expected last month count 1, got %d", last.Count)
	}
	// All Pct values must be in range
	for _, m := range detail.CommitsByMonth {
		if m.Pct < 0 || m.Pct > 100 {
			t.Errorf("Pct out of range [0,100]: %d (month %s)", m.Pct, m.Label)
		}
	}
}

func TestGetRepoDetailCommitsByDay(t *testing.T) {
	// 2024-01-01 = Monday, 2024-01-07 = Sunday
	monday := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	sunday := time.Date(2024, 1, 7, 12, 0, 0, 0, time.UTC)
	r := makeRepoWithCommits(t, []commitSpec{
		{author: "Alice", email: "alice@example.com", message: "mon", when: monday},
		{author: "Alice", email: "alice@example.com", message: "sun", when: sunday},
	})

	detail, err := GetRepoDetail(r, r.Path)
	if err != nil {
		t.Fatalf("GetRepoDetail: %v", err)
	}
	if len(detail.CommitsByDay) != 7 {
		t.Fatalf("expected 7 day entries, got %d", len(detail.CommitsByDay))
	}
	if detail.CommitsByDay[0].Label != "Mon" {
		t.Errorf("expected first entry Mon, got %q", detail.CommitsByDay[0].Label)
	}
	if detail.CommitsByDay[0].Count != 1 {
		t.Errorf("expected Monday count 1, got %d", detail.CommitsByDay[0].Count)
	}
	if detail.CommitsByDay[6].Label != "Sun" {
		t.Errorf("expected last entry Sun, got %q", detail.CommitsByDay[6].Label)
	}
	if detail.CommitsByDay[6].Count != 1 {
		t.Errorf("expected Sunday count 1, got %d", detail.CommitsByDay[6].Count)
	}
}
