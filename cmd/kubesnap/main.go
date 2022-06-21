package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"kubesnap.io/kubesnap/internal/config"
	"kubesnap.io/kubesnap/internal/utilities"
	k8s "kubesnap.io/kubesnap/pkg/kubernetes"
	snap "kubesnap.io/kubesnap/pkg/snap"
)

var (
	// Config vars
	isEventBasedSnaps = config.FetchIsEventBasedSnaps()
	resyncDuration    = config.FetchResyncDuration()
	reasonRegex       = config.FetchReasonRegex()
	isPrintWarnings   = config.FetchIsPrintWarnings()
	lastSeenThreshold = config.FetchLastSeenThreshold()

	// Kuberenetes vars
	clientset          *kubernetes.Clientset
	scheme, serializer = k8s.GenerateSerializer()
)

const (
	// K8s related error messages
	KUBE_CONFIG_FAILED     = "Unable to generate in-cluster config."
	KUBE_CLIENT_GEN_FAILED = "Unable to generate clientset."
	CACHE_SYNC_TIMEOUT     = "timed out waiting for caches to sync"

	// Other constants
	delimiter  = ": "
	WARNING    = "[WARNING]"
	secretsDir = "/etc/secrets/"
)

func main() {

	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	utilities.CheckIfError(err, KUBE_CONFIG_FAILED)

	// Create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	utilities.CheckIfError(err, KUBE_CLIENT_GEN_FAILED)

	// Get resource version for watching events
	factory := informers.NewSharedInformerFactory(clientset, time.Duration(int(math.Pow10(9))*resyncDuration))

	// Get the informer for the right resource
	informer := factory.Core().V1().Events().Informer()

	// Create a channel to stops the shared informer gracefully
	stopper := make(chan struct{})
	defer close(stopper)

	// Kubernetes serves an utility to handle API crashes
	defer runtime.HandleCrash()

	// Trigger custom logic on addition on new event
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    saveCreatedEvent,
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

func saveCreatedEvent(obj interface{}) {
	event := obj.(*corev1.Event)
	message := "Event Creation Detected: " + event.Reason
	saveEvent(message, event.Reason+delimiter+event.Message, event)
}

func saveUpdatedEvent(oldObj interface{}, newObj interface{}) {
	oldEvent := oldObj.(*corev1.Event)
	newEvent := newObj.(*corev1.Event)
	oldMessage := oldEvent.Reason + delimiter + oldEvent.Message
	newMessage := newEvent.Reason + delimiter + newEvent.Message
	message := "Event Update Detected: " + oldEvent.Reason + " => " + newEvent.Reason
	description := oldMessage + " => " + newMessage
	saveEvent(message, description, newEvent)
}

func saveDeletedEvent(obj interface{}) {
	event := obj.(*corev1.Event)
	message := "Event Delete Detected: " + event.Reason
	saveEvent(message, event.Reason+delimiter+event.Message, event)
}

func saveEvent(message string, description string, event *corev1.Event) {
	eventReason := event.Reason
	lastSeenDuration := time.Now().Unix() - event.LastTimestamp.Time.Unix()

	// Print warnings and return
	if isEventBasedSnaps && event.Type == corev1.EventTypeNormal {
		if isPrintWarnings {
			utilities.CreateTimedLog(WARNING, "Ignored normal event:", eventReason)
		}
		return
	} else if !reasonRegex.MatchString(eventReason) {
		if isPrintWarnings {
			utilities.CreateTimedLog(WARNING, "Ignored event:", eventReason, ". Regex match failed.")
		}
		return
	} else if lastSeenDuration > int64(lastSeenThreshold) {
		if isPrintWarnings {
			utilities.CreateTimedLog(WARNING, "Ignored event:", eventReason, ". Last seen duration", fmt.Sprint(lastSeenDuration), "was greater than threshold.")
		}
		return
	}

	// Print settings
	fmt.Println()
	utilities.CreateTimedLog("[INFO] Detected Matching Event:", eventReason, event.Message)
	utilities.CreateTimedLog("[CONFIG] event_based_snaps:", strconv.FormatBool(isEventBasedSnaps))
	utilities.CreateTimedLog("[CONFIG] resync_duration:", fmt.Sprint(resyncDuration))
	utilities.CreateTimedLog("[CONFIG] reason_regex:", reasonRegex.String())
	utilities.CreateTimedLog("[CONFIG] print_warnings:", strconv.FormatBool(isPrintWarnings))

	// Take a snapshot
	snap.TakeSnap(clientset, scheme, serializer, message,
		description, utilities.GetValueOf(secretsDir, "repo-url"),
		utilities.GetValueOf(secretsDir, "repo-branch"))
}
