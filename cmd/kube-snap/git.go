package main

import (
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func CloneRepo() (*git.Repository, error) {
	return git.PlainClone(CloneDir, false, &git.CloneOptions{
		URL:      getValueOf("repo-url", ""),
		Progress: os.Stdout,
	})
}

func AddAll(worktree *git.Worktree) error {
	return worktree.AddWithOptions(&git.AddOptions{
		All:  true,
		Path: CloneDir,
	})
}

func CommitChanges(worktree *git.Worktree) error {
	_, err := worktree.Commit("Snapshot taken on "+time.Now().UTC().String(), &git.CommitOptions{
		Author: &object.Signature{
			Name: getValueOf("commit-author", "kube-snap"),
			When: time.Now(),
		},
	})
	return err
}

func Push(repo *git.Repository) {
	repo.Push(&git.PushOptions{})
}
