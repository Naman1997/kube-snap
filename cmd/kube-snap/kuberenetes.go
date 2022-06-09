package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

func saveNodes(clientset *kubernetes.Clientset, codec runtime.Codec) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	checkIfError(err, "Unable to iterate nodes using the current config.")
	recreateDir(nodePath)

	for index := range nodes.Items {
		node := nodes.Items[index]
		yaml, err := runtime.Encode(codec, &node)
		checkIfError(err, "Unable to encode: "+node.Name+".")
		path := nodePath + node.GetName()
		createFile(path, string(yaml))
	}
}

func saveNamespaces(clientset *kubernetes.Clientset, codec runtime.Codec) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	checkIfError(err, "Unable to iterate namespaces using the current config.")
	recreateDir(namspacePath)

	for index := range namespaces.Items {
		namespace := namespaces.Items[index]
		yaml, err := runtime.Encode(codec, &namespace)
		checkIfError(err, "Unable to encode: "+namespace.Name+".")
		path := namspacePath + namespace.GetName()
		createFile(path, string(yaml))
	}
}
