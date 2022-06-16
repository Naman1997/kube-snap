package main

import (
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
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

func takeSnap(clientset *kubernetes.Clientset, scheme *runtime.Scheme, serializer *json.Serializer, reason string, description string) {

	// Setup workdir
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

func saveKuberentesObjects(clientset *kubernetes.Clientset, scheme *runtime.Scheme, serializer *json.Serializer) {

	codecStorageV1 := k8s.GenerateCodec(scheme, serializer, "storage.k8s.io", "v1")
	codecRbacV1 := k8s.GenerateCodec(scheme, serializer, "rbac.authorization.k8s.io", "v1")
	codecAppsV1 := k8s.GenerateCodec(scheme, serializer, "apps", "v1")
	codecBatchV1 := k8s.GenerateCodec(scheme, serializer, "batch", "v1")
	codecNetworkV1 := k8s.GenerateCodec(scheme, serializer, "networking.k8s.io", "v1")
	codecV1 := k8s.GenerateCodec(scheme, serializer, "", "v1")

	// storage.k8s.io/v1
	saveStorageObjects(clientset, codecStorageV1)

	// rbac.authorization.k8s.io/v1
	saveRbacObjects(clientset, codecRbacV1)

	// v1
	namespaces := saveV1Objects(clientset, codecV1)

	// Namespace Objects
	for _, namespace := range namespaces {

		// v1
		saveNamepacedV1Objects(clientset, codecV1, namespace)

		// apps/v1
		saveNamepacedAppsV1Objects(clientset, codecAppsV1, namespace)

		// batch/v1
		saveNamepacedBatchV1Objects(clientset, codecBatchV1, namespace)

		// networking.k8s.io/v1
		saveNamepacedNetworkingV1Objects(clientset, codecNetworkV1, namespace)

		// rbac.authorization.k8s.io/v1
		saveNamepacedRbacV1Objects(clientset, codecRbacV1, namespace)
	}
}

func saveStorageObjects(clientset *kubernetes.Clientset, codecStorageV1 runtime.Codec) {
	k8s.SaveStorageClasses(clientset, codecStorageV1)
}

func saveRbacObjects(clientset *kubernetes.Clientset, codecRbacV1 runtime.Codec) {
	k8s.SaveClusterRoleBindings(clientset, codecRbacV1)
	k8s.SaveClusterRoles(clientset, codecRbacV1)
}

func saveV1Objects(clientset *kubernetes.Clientset, codecV1 runtime.Codec) []string {
	k8s.SavePersistentVolumes(clientset, codecV1)
	k8s.SaveNodes(clientset, codecV1)
	namespaces := k8s.SaveNamespaces(clientset, codecV1)
	return namespaces
}

func saveNamepacedV1Objects(clientset *kubernetes.Clientset, codecV1 runtime.Codec, namespace string) {
	k8s.SaveConfigMaps(clientset, codecV1, namespace)
	k8s.SavePersistentVolumeClaims(clientset, codecV1, namespace)
	k8s.SavePods(clientset, codecV1, namespace)
	k8s.SaveSecrets(clientset, codecV1, namespace)
	k8s.SaveServiceAccounts(clientset, codecV1, namespace)
	k8s.SaveServices(clientset, codecV1, namespace)
}

func saveNamepacedAppsV1Objects(clientset *kubernetes.Clientset, codecAppsV1 runtime.Codec, namespace string) {
	k8s.SaveDaemonsets(clientset, codecAppsV1, namespace)
	k8s.SaveDeployments(clientset, codecAppsV1, namespace)
	k8s.SaveReplicaSets(clientset, codecAppsV1, namespace)
	k8s.SaveStatefulSets(clientset, codecAppsV1, namespace)
}

func saveNamepacedBatchV1Objects(clientset *kubernetes.Clientset, codecBatchV1 runtime.Codec, namespace string) {
	k8s.SaveCronJobs(clientset, codecBatchV1, namespace)
	k8s.SaveJobs(clientset, codecBatchV1, namespace)
}

func saveNamepacedNetworkingV1Objects(clientset *kubernetes.Clientset, codecNetworkV1 runtime.Codec, namespace string) {
	k8s.SaveIngresses(clientset, codecNetworkV1, namespace)
}

func saveNamepacedRbacV1Objects(clientset *kubernetes.Clientset, codecRbacV1 runtime.Codec, namespace string) {
	k8s.SaveRoleBindings(clientset, codecRbacV1, namespace)
	k8s.SaveRoles(clientset, codecRbacV1, namespace)
}
