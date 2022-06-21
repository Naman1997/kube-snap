package snap

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	"kubesnap.io/kubesnap/internal/git"
	"kubesnap.io/kubesnap/internal/utilities"
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
	GIT_PUSH_COMPLETED = INFO + "Created new snap!"

	// Other constants
	cloneDir    = "/repo"
	gitSubDir   = ".git"
	authorName  = "kubesnap"
	authorEmail = "kubesnap@kubesnap.com"
)

func TakeSnap(clientset *kubernetes.Clientset, scheme *runtime.Scheme, serializer *json.Serializer,
	reason string, description string, repoUrl string, branch string) {

	// Setup workdir
	// Make sure clone dir exists
	_, err := os.Stat(cloneDir)
	utilities.CheckIfError(err, CLONE_DIR_NOT_DETECTED)

	// Change dir to repo
	err = os.Chdir(cloneDir)
	utilities.CheckIfError(err, CHDIR_FAILED)

	_, err = os.Stat(gitSubDir)
	if os.IsNotExist(err) { // Repo is not cloned yet. Start cloning.
		utilities.CheckIfError(git.CloneRepo(repoUrl), GIT_CLONE_FAILED)
		utilities.CreateTimedLog(GIT_CLONE_PASSED)
	} else { // Repo is already cloned. Pull latest changes
		utilities.CheckIfError(git.PullOrigin(), GIT_PULL_FAILED)
		utilities.CreateTimedLog(GIT_PULL_PASSED)
	}

	// Configure git author
	git.SetupAuthor(authorEmail, authorName)

	// Checkout in branch
	if len(branch) > 0 {
		utilities.CreateTimedLog(GIT_SWITCH_PASSED)
		utilities.CheckIfError(git.SwitchBranch(branch), GIT_SWITCH_FAILED)
	}

	// Save k8s objects
	saveKuberentesObjects(clientset, scheme, serializer)

	// Add all files
	err = git.AddAll()
	utilities.CheckIfError(err, GIT_ADD_ALL_FAILED)
	utilities.CreateTimedLog(GIT_ADD_ALL_PASSED)

	// Check status of changes
	o, err := git.Status()
	utilities.CheckIfError(err, GIT_STATUS_FAILED)
	utilities.CreateTimedLog(GIT_STATUS_PASSED)
	if len(o) > 0 {
		fmt.Print("["+time.Now().UTC().Format(time.UnixDate)+"]", INFO+o)
		// Commit all files
		utilities.CheckIfError(git.CommitChanges(reason, description), GIT_COMMIT_FAILED)
		utilities.CreateTimedLog(GIT_COMMIT_PASSED)

		// Git push using default options
		utilities.CheckIfError(git.Push(), GIT_PUSH_FAILED)
		utilities.CreateTimedLog(GIT_PUSH_PASSED)
		utilities.CreateTimedLog(INFO + GIT_PUSH_COMPLETED)
	} else {
		// No changes
		utilities.CreateTimedLog(GIT_NO_DIFF_FOUND)
	}

}
