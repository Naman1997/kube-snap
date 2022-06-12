package main

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"kube-snap.io/kube-snap/internal/git"
	"kube-snap.io/kube-snap/internal/utilities"
	k8s "kube-snap.io/kube-snap/pkg/kubernetes"
)

const (
	// Filesystem related error messages
	CLONE_DIR_NOT_DETECTED = "Clone dir not detected. [PV not mounted?]"
	CHDIR_FAILED           = "Unable to change dir to the cloned repo."

	// Git related error messages
	GIT_CLONE_FAILED   = "Unable to clone repo."
	GIT_PULL_FAILED    = "Unable to execute git pull."
	GIT_SWITCH_FAILED  = "Unable to switch branch."
	GIT_ADD_ALL_FAILED = "Unable to execute git add --all."
	GIT_STATUS_FAILED  = "Unable to execute git status --porcelain."
	GIT_COMMIT_FAILED  = "Unable to execute git commit."
	GIT_PUSH_FAILED    = "Unable to execute git push."

	// Git related info messages
	INFO               = "[INFO] "
	GIT_CLONE_PASSED   = INFO + "Executing git clone."
	GIT_PULL_PASSED    = INFO + "Executing git pull."
	GIT_SWITCH_PASSED  = INFO + "Executing git switch."
	GIT_ADD_ALL_PASSED = INFO + "Executing git add --all."
	GIT_STATUS_PASSED  = INFO + "Executing git status --porcelain."
	GIT_COMMIT_PASSED  = INFO + "Executing git commit."
	GIT_PUSH_PASSED    = INFO + "Executing git push."
	GIT_NO_DIFF_FOUND  = INFO + "No diff found! Branch is up to date."

	// Other constants
	cloneDir    = "/repo"
	gitSubDir   = ".git"
	authorName  = "kube-snap"
	authorEmail = "kubesnap@kubesnap.com"
	secretsDir  = "/etc/secrets/"
)

func TakeSnap(clientset *kubernetes.Clientset, codec runtime.Codec, reason string, description string) {

	// Setup workdir
	setupGitRepo()

	// Save k8s objects
	saveKuberentesObjects(clientset, codec)

	// Add all files
	err := git.AddAll()
	utilities.CheckIfError(err, GIT_ADD_ALL_FAILED)
	utilities.CreateTimedLog(GIT_ADD_ALL_PASSED)

	// Check status of changes
	o, err := git.Status()
	utilities.CheckIfError(err, GIT_STATUS_FAILED)
	utilities.CreateTimedLog(GIT_STATUS_PASSED)
	if len(o) > 0 {
		fmt.Print("["+time.Now().UTC().String()+"]", INFO+o)
		// Commit all files
		utilities.CheckIfError(git.CommitChanges(reason, description), GIT_COMMIT_FAILED)
		utilities.CreateTimedLog(GIT_COMMIT_PASSED)

		// Git push using default options
		utilities.CheckIfError(git.Push(), GIT_PUSH_FAILED)
		utilities.CreateTimedLog(GIT_PUSH_PASSED)

	} else {
		// No changes
		utilities.CreateTimedLog(GIT_NO_DIFF_FOUND)
	}

}

func setupGitRepo() {
	// Make sure clone dir exists
	_, err := os.Stat(cloneDir)
	utilities.CheckIfError(err, CLONE_DIR_NOT_DETECTED)

	// Change dir to repo
	err = os.Chdir(cloneDir)
	utilities.CheckIfError(err, CHDIR_FAILED)

	_, err = os.Stat(gitSubDir)
	if os.IsNotExist(err) { // Repo is not cloned yet. Start cloning.
		utilities.CheckIfError(git.CloneRepo(utilities.GetValueOf(secretsDir, "repo-url")), GIT_CLONE_FAILED)
		utilities.CreateTimedLog(GIT_CLONE_PASSED)
	} else { // Repo is already cloned. Pull latest changes
		utilities.CheckIfError(git.PullOrigin(), GIT_PULL_FAILED)
		utilities.CreateTimedLog(GIT_PULL_PASSED)
	}

	// Configure git author
	git.SetupAuthor(authorEmail, authorName)

	// Checkout in branch
	branch := utilities.GetValueOf(secretsDir, "repo-branch")
	if len(branch) > 0 {
		utilities.CreateTimedLog(GIT_SWITCH_PASSED)
		utilities.CheckIfError(git.SwitchBranch(branch), GIT_SWITCH_FAILED)
	}
}

func saveKuberentesObjects(clientset *kubernetes.Clientset, codec runtime.Codec) {
	k8s.SaveNodes(clientset, codec)
	k8s.SaveNamespaces(clientset, codec)
}
