package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

func generateCodec() runtime.Codec {
	scheme := scheme.Scheme
	addToScheme(scheme)
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
	return versioning.NewDefaultingCodecForScheme(
		scheme,
		serializer,
		serializer,
		schema.GroupVersion{Version: Version},
		runtime.InternalGroupVersioner,
	)
}

func saveNodes(clientset *kubernetes.Clientset, codec runtime.Codec) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	checkIfError(err)
	createDir(nodePath)
	existing := findExistingConfigs(nodePath)
	var current []string

	for index := range nodes.Items {
		node := nodes.Items[index]
		yaml, err := runtime.Encode(codec, &node)
		checkIfError(err)
		path := nodePath + node.GetName()
		current = append(current, path)
		createFile(path, string(yaml))
	}
	deleteMissingConfigs(current, existing)
}

func saveNamespaces(clientset *kubernetes.Clientset, codec runtime.Codec) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	checkIfError(err)
	createDir(namspacePath)
	existing := findExistingConfigs(namspacePath)
	var current []string

	for index := range namespaces.Items {
		namespace := namespaces.Items[index]
		yaml, err := runtime.Encode(codec, &namespace)
		checkIfError(err)
		path := namspacePath + namespace.GetName()
		current = append(current, path)
		createFile(path, string(yaml))
	}
	deleteMissingConfigs(current, existing)
}
