package snap

import (
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes"
	k8s "kubesnap.io/kubesnap/pkg/kubernetes"
)

func saveKuberentesObjects(clientset *kubernetes.Clientset, scheme *runtime.Scheme, serializer *json.Serializer) {

	// Create codecs for different groups and versions
	codecStorageV1 := k8s.GenerateCodec(scheme, serializer, "storage.k8s.io", "v1")
	codecRbacV1 := k8s.GenerateCodec(scheme, serializer, "rbac.authorization.k8s.io", "v1")
	codecAppsV1 := k8s.GenerateCodec(scheme, serializer, "apps", "v1")
	codecBatchV1 := k8s.GenerateCodec(scheme, serializer, "batch", "v1")
	codecNetworkV1 := k8s.GenerateCodec(scheme, serializer, "networking.k8s.io", "v1")
	codecV1 := k8s.GenerateCodec(scheme, serializer, "", "v1")

	// Create a wait group for non-namespaced methods
	wg := new(sync.WaitGroup)
	wg.Add(5)

	// v1
	namespaces := k8s.SaveNamespaces(clientset, codecV1)
	go k8s.SavePersistentVolumes(clientset, codecV1, wg)
	go k8s.SaveNodes(clientset, codecV1, wg)

	// storage.k8s.io/v1
	go k8s.SaveStorageClasses(clientset, codecStorageV1, wg)

	// rbac.authorization.k8s.io/v1
	go k8s.SaveClusterRoleBindings(clientset, codecRbacV1, wg)
	go k8s.SaveClusterRoles(clientset, codecRbacV1, wg)

	// Wait for non-namespaced methods to finish
	wg.Wait()

	// Namespace Objects
	for _, namespace := range namespaces {

		/*
			Not creating wait groups for namespaced resources
			as they tend to cause client-side throttling
			https://kubernetes.io/docs/concepts/cluster-administration/flow-control/
		*/

		// v1
		k8s.SaveConfigMaps(clientset, codecV1, namespace)
		k8s.SavePersistentVolumeClaims(clientset, codecV1, namespace)
		k8s.SavePods(clientset, codecV1, namespace)
		k8s.SaveSecrets(clientset, codecV1, namespace)
		k8s.SaveServiceAccounts(clientset, codecV1, namespace)
		k8s.SaveServices(clientset, codecV1, namespace)

		// apps/v1
		k8s.SaveDaemonsets(clientset, codecAppsV1, namespace)
		k8s.SaveDeployments(clientset, codecAppsV1, namespace)
		k8s.SaveReplicaSets(clientset, codecAppsV1, namespace)
		k8s.SaveStatefulSets(clientset, codecAppsV1, namespace)

		// batch/v1
		k8s.SaveCronJobs(clientset, codecBatchV1, namespace)
		k8s.SaveJobs(clientset, codecBatchV1, namespace)

		// networking.k8s.io/v1
		k8s.SaveIngresses(clientset, codecNetworkV1, namespace)

		// rbac.authorization.k8s.io/v1
		k8s.SaveRoleBindings(clientset, codecRbacV1, namespace)
		k8s.SaveRoles(clientset, codecRbacV1, namespace)
	}
}
