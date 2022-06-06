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

	// Git clone a repo using creds from a secret
	repo, err := cloneRepo()
	fmt.Println()
	if err != nil {
		fmt.Println("Git Clone Failed!")
		checkIfError(err)
	}

	// Generate worktree
	worktree, err := repo.Worktree()
	checkIfError(err)

	// TODO: Figure out how to switch branches using git-go
	// https://github.com/go-git/go-git/issues/241

	// Checkout provided branch
	// fmt.Println("Switching to: " + plumbing.NewBranchReferenceName(getValueOf("repo-branch", "")))
	// err = worktree.Checkout(&git.CheckoutOptions{
	// 	Create: true,
	// 	Branch: plumbing.NewBranchReferenceName(getValueOf("repo-branch", "")),
	// })
	// checkIfError(err)

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	checkIfError(err)

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	checkIfError(err)

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
	checkIfError(err)

	// Check if worktree is clean
	fmt.Println("Checking git status.")
	if checkCleanStatus(worktree) {
		fmt.Println("Clean worktree. Nothing to commit.")
	} else {
		// Commit all files
		fmt.Println("Executing git commit.")
		err = commitChanges(worktree)
		checkIfError(err)

		// Git push using default options
		fmt.Println("Executing git push.")
		push(repo)
	}
}
