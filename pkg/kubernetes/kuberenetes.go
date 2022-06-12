package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"kube-snap.io/kube-snap/internal/utilities"
)

const (
	nodePath     = "nodes/"
	namspacePath = "namespaces/"
)

func SaveNodes(clientset *kubernetes.Clientset, codec runtime.Codec) {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate nodes using the current config.")
	utilities.RecreateDir(nodePath)

	for index := range nodes.Items {
		node := nodes.Items[index]
		yaml, err := runtime.Encode(codec, &node)
		utilities.CheckIfError(err, "Unable to encode: "+node.Name+".")
		path := nodePath + node.GetName()
		utilities.CreateFile(path, string(yaml))
	}
}

func SaveNamespaces(clientset *kubernetes.Clientset, codec runtime.Codec) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate namespaces using the current config.")
	utilities.RecreateDir(namspacePath)

	for index := range namespaces.Items {
		namespace := namespaces.Items[index]
		yaml, err := runtime.Encode(codec, &namespace)
		utilities.CheckIfError(err, "Unable to encode: "+namespace.Name+".")
		path := namspacePath + namespace.GetName()
		utilities.CreateFile(path, string(yaml))
	}
}
