package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kube-snap.io/kube-snap/internal/utilities"
	k8s "kube-snap.io/kube-snap/pkg/kubernetes"
)

func main() {

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	utilities.CheckIfError(err, K8S_CONFIG_FAILED)

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	utilities.CheckIfError(err, K8s_CLIENT_GEN_FAILED)

	// Generate codec for serialization
	codec := k8s.GenerateCodec()

	TakeSnap(clientset, codec)
}
