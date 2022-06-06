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

	// Create the nodes dir
	createDir(nodePath)

	for node_index := range nodes.Items {
		node := nodes.Items[node_index]
		yaml, err := runtime.Encode(codec, &node)
		checkIfError(err)
		createFile(nodePath+node.GetName(), string(yaml))
	}
}

func saveNamespaces(clientset *kubernetes.Clientset, codec runtime.Codec) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	checkIfError(err)

	// Create the nodes dir
	createDir(namspacePath)

	for node_index := range namespaces.Items {
		namespace := namespaces.Items[node_index]
		yaml, err := runtime.Encode(codec, &namespace)
		checkIfError(err)
		createFile(namspacePath+namespace.GetName(), string(yaml))
	}
}
