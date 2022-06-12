package main

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"kube-snap.io/kube-snap/internal/utilities"
	k8s "kube-snap.io/kube-snap/pkg/kubernetes"
)

var (
	clientset *kubernetes.Clientset
	codec     = k8s.GenerateCodec()
)

func main() {

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	utilities.CheckIfError(err, KUBE_CONFIG_FAILED)

	// Create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	utilities.CheckIfError(err, KUBE_CLIENT_GEN_FAILED)

	// Get resource version for watching events
	factory := informers.NewSharedInformerFactory(clientset, 0)

	// Get the informer for the right resource
	informer := factory.Core().V1().Events().Informer()

	// Create a channel to stops the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()

	// Trigger custom logic on addition on new event
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    saveEvent,
		UpdateFunc: saveUpdatedEvent,
		DeleteFunc: saveDeletedEvent,
	})

	go informer.Run(stopper)
	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		runtime.HandleError(fmt.Errorf(CACHE_SYNC_TIMEOUT))
		return
	}
	<-stopper
}

func saveEvent(obj interface{}) {
	// Cast the obj as event
	event := obj.(*corev1.Event)
	message := event.Reason + delimiter + event.Message
	TakeSnap(clientset, codec, message)
}

func saveUpdatedEvent(oldObj interface{}, newObj interface{}) {
	// Cast the obj as event
	oldEvent := oldObj.(*corev1.Event)
	newEvent := newObj.(*corev1.Event)
	oldMessage := oldEvent.Reason + delimiter + oldEvent.Message
	newReason := newEvent.Reason + delimiter + newEvent.Message
	reason := "Event Update Detected: " + oldMessage + " to " + newReason
	TakeSnap(clientset, codec, reason)
}

func saveDeletedEvent(obj interface{}) {
	// Cast the obj as event
	event := obj.(*corev1.Event)
	message := "Event Delete Detected: " + event.Reason + delimiter + event.Message
	TakeSnap(clientset, codec, message)
}
