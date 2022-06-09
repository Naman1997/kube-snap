package main

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	Version      = "v1"
	CloneDir     = "repo"
	nodePath     = "nodes/"
	namspacePath = "namespaces/"
)

func main() {

	// Create the CloneDir
	createDir(CloneDir)

	// TODO: Only pull the latest changes if repo exists on disk
	// Clone the repo
	fmt.Println("Executing git clone.")
	err := cloneRepo()
	if err != nil {
		checkIfError(err, "Unable to clone repo.")
	}

	// Change dir to repo
	err = os.Chdir(CloneDir)
	checkIfError(err, "Unable to change dir to the cloned repo.")

	// Configure git author
	setupAuthor()

	// Checkout in branch
	err = switchBranch()
	checkIfError(err, "Unable to switch branch")

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	checkIfError(err, "Unable to generate in-cluster config.")

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	checkIfError(err, "Unable to generate clientset.")

	// Generate codec for serialization
	codec := generateCodec()

	//Save all nodes
	saveNodes(clientset, codec)

	//Save all namespaces
	saveNamespaces(clientset, codec)

	// Add all files
	fmt.Println()
	fmt.Println("Executing git add.")
	err = addAll()
	checkIfError(err, "Unable to execute git add.")

	// Check status of changes
	fmt.Println("Executing git status --porcelain")
	o, err := status()
	checkIfError(err, "Unable to execute git status.")
	fmt.Println(o)
	if len(o) > 0 {
		// Commit all files
		fmt.Println("Executing git commit.")
		err = commitChanges()
		checkIfError(err, "Unable to execute git commit.")

		// Git push using default options
		fmt.Println("Executing git push.")
		err = push()
		checkIfError(err, "Unable to execute git push.")

	} else {
		// No changes
		fmt.Println("No diff found! Branch", getValueOf("repo-branch"), "is up to date.")
	}

}
