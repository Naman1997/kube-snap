package main

import (
	"context"
	"fmt"
	"path/filepath"

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

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const Version = "v1"

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

	/*
		TODO:
		- Write files in a specific directory structure
		- Git push to the same specified branch
	*/

	// Git clone a repo using creds from a secret
	_, err := git.PlainClone("./repo/", false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: getValueOf("repo-user", "abc123"),
			Password: getValueOf("repo-pass", ""),
		},
		URL:        getValueOf("repo-url", ""),
		RemoteName: getValueOf("repo-branch", ""),
		Progress:   os.Stdout,
	})

	// Remove all non dot files and dirs in the cloned dir
	files, err := filepath.Glob("./repo/.*")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}

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
		yaml, err := runtime.Encode(Codec, &nodes.Items[node_index])
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%+v\n", string(yaml))
		fmt.Println()
	}

	//Get all namespaces
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for node_index := range namespaces.Items {
		yaml, err := runtime.Encode(Codec, &namespaces.Items[node_index])
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%+v\n", string(yaml))
		fmt.Println()
	}
}
