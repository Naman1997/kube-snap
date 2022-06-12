package main

const (
	// Common constants
	cloneDir    = "repo"
	gitSubDir   = ".git"
	authorName  = "kube-snap"
	authorEmail = "kubesnap@kubesnap.com"

	// K8s related error messages
	K8S_CONFIG_FAILED     = "Unable to generate in-cluster config."
	K8s_CLIENT_GEN_FAILED = "Unable to generate clientset."

	// Filesystem related error messages
	CLONE_DIR_NOT_DETECTED = "Clone dir not detected. PV not mounted?"
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
)
