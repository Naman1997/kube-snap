package main

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	Version      = "v1"
	CloneDir     = "/repo"
	nodePath     = CloneDir + "/nodes/"
	namspacePath = CloneDir + "/namespaces/"
)

func main() {

	// Create the CloneDir
	createDir(CloneDir)

	// TODO: Only pull the latest changes if repo exists on disk
	// Clone the repo
	repo, err := cloneRepo()
	fmt.Println()
	if err != nil {
		checkIfError(err, "Unable to clone repo")
	}

	// Generate worktree
	worktree, err := repo.Worktree()
	checkIfError(err, "Unable to generate worktree")

	// TODO: Figure out how to switch branches using git-go
	// https://github.com/go-git/go-git/issues/241

	// Checkout in branch
	// err = checkout(worktree)
	// checkIfError(err, "Unable to checkout in branch")

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	checkIfError(err, "Unable to generate in-cluster config")

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	checkIfError(err, "Unable to generate clientset")

	// Generate codec for serialization
	codec := generateCodec()

	//Save all nodes
	saveNodes(clientset, codec)

	//Save all namespaces
	saveNamespaces(clientset, codec)

	// Add all files
	fmt.Println()
	fmt.Println("Executing git add --all.")
	err = addAll(worktree)
	checkIfError(err, "Unable to execute git add --all")

	// Check if worktree is clean
	fmt.Println("Checking git status.")
	if checkCleanStatus(worktree) {
		fmt.Println("Clean worktree. Nothing to commit.")
	} else {
		// Commit all files
		fmt.Println("Executing git commit.")
		_, err := commitChanges(worktree)
		checkIfError(err, "Unable to execute git commit")

		// Git push using default options
		fmt.Println("Executing git push.")
		push(repo)
	}
}
