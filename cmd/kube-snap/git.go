package main

import (
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func cloneRepo() (*git.Repository, error) {
	return git.PlainClone(CloneDir, false, &git.CloneOptions{
		URL:      getValueOf("repo-url", ""),
		Progress: os.Stdout,
	})
}

// func checkout(worktree *git.Worktree) error {
// 	return worktree.Checkout(&git.CheckoutOptions{
// 		Hash:   plumbing.Hash{},
// 		Branch: plumbing.NewBranchReferenceName(getValueOf("repo-branch", "")),
// 		Force:  true,
// 	})
// }

func addAll(worktree *git.Worktree) error {
	// TODO: Fix add all [Maybe maually send the list of files changed]
	return worktree.AddWithOptions(&git.AddOptions{
		All:  true,
		Glob: CloneDir,
	})
}

func checkCleanStatus(worktree *git.Worktree) bool {
	status, err := worktree.Status()
	checkIfError(err, "Unable to find worktree status.")
	return status.IsClean()
}

func commitChanges(worktree *git.Worktree) (plumbing.Hash, error) {
	commit, err := worktree.Commit("Snapshot taken on "+time.Now().UTC().String(), &git.CommitOptions{
		Author: &object.Signature{
			Name: "kube-snap",
			When: time.Now(),
		},
	})
	return commit, err
}

func push(repo *git.Repository) {
	repo.Push(&git.PushOptions{
		Force: true,
	})
}
