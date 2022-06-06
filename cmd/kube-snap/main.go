package main

import (
	"context"
	"fmt"

	"os"
	// "path/filepath"

	appsv1 "github.com/openshift/api/apps/v1"
	authorizationv1 "github.com/openshift/api/authorization/v1"
	buildv1 "github.com/openshift/api/build/v1"
	imagev1 "github.com/openshift/api/image/v1"
	networkv1 "github.com/openshift/api/network/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	projectv1 "github.com/openshift/api/project/v1"
	quotav1 "github.com/openshift/api/quota/v1"
	routev1 "github.com/openshift/api/route/v1"
	securityv1 "github.com/openshift/api/security/v1"
	templatev1 "github.com/openshift/api/template/v1"
	userv1 "github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const (
	Version  = "v1"
	CloneDir = "/repo"
)

var Codec runtime.Codec

func init() {
	appsv1.AddToScheme(scheme.Scheme)
	authorizationv1.AddToScheme(scheme.Scheme)
	buildv1.AddToScheme(scheme.Scheme)
	imagev1.AddToScheme(scheme.Scheme)
	networkv1.AddToScheme(scheme.Scheme)
	oauthv1.AddToScheme(scheme.Scheme)
	projectv1.AddToScheme(scheme.Scheme)
	quotav1.AddToScheme(scheme.Scheme)
	routev1.AddToScheme(scheme.Scheme)
	securityv1.AddToScheme(scheme.Scheme)
	templatev1.AddToScheme(scheme.Scheme)
	userv1.AddToScheme(scheme.Scheme)
	appsv1.AddToSchemeInCoreGroup(scheme.Scheme)
	authorizationv1.AddToSchemeInCoreGroup(scheme.Scheme)
	buildv1.AddToSchemeInCoreGroup(scheme.Scheme)
	imagev1.AddToSchemeInCoreGroup(scheme.Scheme)
	networkv1.AddToSchemeInCoreGroup(scheme.Scheme)
	oauthv1.AddToSchemeInCoreGroup(scheme.Scheme)
	projectv1.AddToSchemeInCoreGroup(scheme.Scheme)
	quotav1.AddToSchemeInCoreGroup(scheme.Scheme)
	routev1.AddToSchemeInCoreGroup(scheme.Scheme)
	securityv1.AddToSchemeInCoreGroup(scheme.Scheme)
	templatev1.AddToSchemeInCoreGroup(scheme.Scheme)
	userv1.AddToSchemeInCoreGroup(scheme.Scheme)
}

func main() {

	// Create the CloneDir
	if err := os.Mkdir(CloneDir, os.ModePerm); err != nil {
		panic(err)
	}

	// Git clone a repo using creds from a secret
	repo, err := CloneRepo()
	if err != nil {
		fmt.Println("Git Clone Failed!")
		panic(err)
	}

	// Generate worktree
	worktree, err := repo.Worktree()
	CheckIfError(err)

	// TODO: Figure out how to switch branches using git-go
	// https://github.com/go-git/go-git/issues/241

	// Checkout provided branch
	// fmt.Println("Switching to: " + plumbing.NewBranchReferenceName(getValueOf("repo-branch", "")))
	// err = worktree.Checkout(&git.CheckoutOptions{
	// 	Create: true,
	// 	Branch: plumbing.NewBranchReferenceName(getValueOf("repo-branch", "")),
	// })
	// CheckIfError(err)

	s := scheme.Scheme
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, s, s)
	Codec = versioning.NewDefaultingCodecForScheme(
		s,
		serializer,
		serializer,
		schema.GroupVersion{Version: Version},
		runtime.InternalGroupVersioner,
	)

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//Get all nodes
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for node_index := range nodes.Items {
		node := nodes.Items[node_index]
		yaml, err := runtime.Encode(Codec, &node)
		if err != nil {
			panic(err.Error())
		}
		createFile(CloneDir+"/"+node.GetName(), string(yaml))
		fmt.Println()
	}

	//Get all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for node_index := range namespaces.Items {
		namespace := namespaces.Items[node_index]
		yaml, err := runtime.Encode(Codec, &namespace)
		if err != nil {
			panic(err.Error())
		}
		createFile(CloneDir+"/"+namespace.GetName(), string(yaml))
		fmt.Println()
	}

	// Add all files
	fmt.Println("Executing git add --all.")
	err = AddAll(worktree)
	CheckIfError(err)

	// Commit all files
	fmt.Println("Executing git commit.")
	err = CommitChanges(worktree)
	CheckIfError(err)

	// Git push using default options
	fmt.Println("Executing git push.")
	Push(repo)
}
