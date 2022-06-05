package main

import (
	"context"
	"fmt"

	// "os"
	// "path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"

	// git "github.com/go-git/go-git/v5"
	// "github.com/go-git/go-git/v5/plumbing/transport/http"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	/*
		TODO:
		- Git clone a repo using creds from a scret
		- Remove all non dot files and dirs in the cloned dir
		- Write files in a specific directory structure
		- Git push to the same specified branch
	*/

	// _, err := git.PlainClone("/repo/", false, &git.CloneOptions{
	// 	Auth: &http.BasicAuth{
	// 		Username: getValueOf("repo-user", "abc123"),
	// 		Password: getValueOf("repo-pass", ""),
	// 	},
	// 	URL:        getValueOf("repo-url", ""),
	// 	RemoteName: getValueOf("repo-branch", ""),
	// 	Progress:   os.Stdout,
	// })

	// files, err := filepath.Glob("/repo/.*")
	// if err != nil {
	// 	panic(err)
	// }
	// for _, f := range files {
	// 	if err := os.Remove(f); err != nil {
	// 		panic(err)
	// 	}
	// }

	// Serializer = Decoder + Encoder.
	serializer := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory,
		scheme.Scheme,
		scheme.Scheme,
		jsonserializer.SerializerOptions{
			Yaml:   true,
			Pretty: true,
			Strict: false,
		},
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
		yaml, err := runtime.Encode(serializer, &nodes.Items[node_index])
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
		yaml, err := runtime.Encode(serializer, &namespaces.Items[node_index])
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%+v\n", string(yaml))
		fmt.Println()
	}
}
