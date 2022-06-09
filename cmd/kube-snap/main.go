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
	gitSubDir    = ".git"
)

func main() {

	// Gather stats of clone dir
	_, err := os.Stat(CloneDir)
	checkIfError(err, "Clone dir not detected. PV not mounted?")

	// Change dir to repo
	err = os.Chdir(CloneDir)
	checkIfError(err, "Unable to change dir to the cloned repo.")

	// Gather stats of git sub dir
	_, err = os.Stat(gitSubDir)

	if os.IsNotExist(err) {
		// Clone the repo
		fmt.Println("[INFO] Executing git clone.")
		checkIfError(cloneRepo(), "Unable to clone repo.")
	} else {
		// Execute git pull
		fmt.Println("[INFO] Executing git pull.")
		checkIfError(pullOrigin(), "Unable to execute git pull.")
	}

	// Configure git author
	setupAuthor()

	// Checkout in branch
	branch := getValueOf("repo-branch")
	if len(branch) > 0 {
		fmt.Println("[INFO] Executing git switch.")
		checkIfError(switchBranch(branch), "Unable to switch branch")
	}

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
	fmt.Println("[INFO] Executing git add.")
	err = addAll()
	checkIfError(err, "Unable to execute git add.")

	// Check status of changes
	fmt.Println("[INFO] Executing git status --porcelain")
	o, err := status()
	checkIfError(err, "Unable to execute git status.")
	fmt.Print(o)
	if len(o) > 0 {
		// Commit all files
		fmt.Println("[INFO] Executing git commit.")
		checkIfError(commitChanges(), "Unable to execute git commit.")

		// Git push using default options
		fmt.Println("[INFO] Executing git push.")
		checkIfError(push(), "Unable to execute git push.")

	} else {
		// No changes
		fmt.Println("[INFO] No diff found! Branch", getValueOf("repo-branch"), "is up to date.")
	}

}
