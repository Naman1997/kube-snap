package main

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"kube-snap.io/kube-snap/internal/git"
	"kube-snap.io/kube-snap/internal/utilities"
	k8s "kube-snap.io/kube-snap/pkg/kubernetes"
)

func TakeSnap(clientset *kubernetes.Clientset, codec runtime.Codec, reason string) {

	// Setup workdir
	setupGitRepo()

	// Save k8s objects
	saveKuberentesObjects(clientset, codec)

	// Add all files
	err := git.AddAll()
	utilities.CheckIfError(err, GIT_ADD_ALL_FAILED)
	fmt.Println(GIT_ADD_ALL_PASSED)

	// Check status of changes
	o, err := git.Status()
	utilities.CheckIfError(err, GIT_STATUS_FAILED)
	fmt.Println(GIT_STATUS_PASSED)
	fmt.Print(o)
	if len(o) > 0 {
		// Commit all files
		utilities.CheckIfError(git.CommitChanges(reason), GIT_COMMIT_FAILED)
		fmt.Println(GIT_COMMIT_PASSED)

		// Git push using default options
		utilities.CheckIfError(git.Push(), GIT_PUSH_FAILED)
		fmt.Println(GIT_PUSH_PASSED)

	} else {
		// No changes
		fmt.Println(GIT_NO_DIFF_FOUND)
	}

}

func setupGitRepo() {
	// Gather stats of clone dir
	_, err := os.Stat(cloneDir)
	utilities.CheckIfError(err, CLONE_DIR_NOT_DETECTED)

	// Change dir to repo
	err = os.Chdir(cloneDir)
	utilities.CheckIfError(err, CHDIR_FAILED)

	_, err = os.Stat(gitSubDir)
	if os.IsNotExist(err) { // Repo is not cloned yet. Start cloning.
		utilities.CheckIfError(git.CloneRepo(utilities.GetValueOf("repo-url")), GIT_CLONE_FAILED)
		fmt.Println(GIT_CLONE_PASSED)
	} else { // Repo is already cloned. Pull latest changes
		utilities.CheckIfError(git.PullOrigin(), GIT_PULL_FAILED)
		fmt.Println(GIT_PULL_PASSED)
	}

	// Configure git author
	git.SetupAuthor(authorEmail, authorName)

	// Checkout in branch
	branch := utilities.GetValueOf("repo-branch")
	if len(branch) > 0 {
		fmt.Println(GIT_SWITCH_PASSED)
		utilities.CheckIfError(git.SwitchBranch(branch), GIT_SWITCH_FAILED)
	}
}

func saveKuberentesObjects(clientset *kubernetes.Clientset, codec runtime.Codec) {
	k8s.SaveNodes(clientset, codec)
	k8s.SaveNamespaces(clientset, codec)
}
