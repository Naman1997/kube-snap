package kubernetes

import (
	"context"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"kubesnap.io/kubesnap/internal/utilities"
)

const (
	nodePath                  = "nodes/"
	storageclassPath          = "storageclasses/"
	clusterrolebindingPath    = "clusterrolebindings/"
	clusterrolePath           = "clusterroles/"
	persistentVolumePath      = "persistent_volumes/"
	namspacePath              = "namespaces/"
	configMapPath             = "configmaps/"
	persistentVolumeClaimPath = "persistent_volume_claims/"
	podPath                   = "pods/"
	secretPath                = "secrets/"
	serviceAccountPath        = "service_accounts/"
	servicePath               = "services/"
	daemonsetPath             = "daemonsets/"
	deploymentPath            = "deployments/"
	replicasetPath            = "replicasets/"
	statefulsetPath           = "statefulsets/"
	cronjobPath               = "cronjobs/"
	jobPath                   = "jobs/"
	ingressPath               = "ingresses/"
	rolebindingPath           = "rolebindings/"
	rolePath                  = "roles/"
)

func SaveNodes(clientset *kubernetes.Clientset, codec runtime.Codec, wg *sync.WaitGroup) {
	defer wg.Done()
	kubeObject, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate nodes using the current config.")
	utilities.RecreateDir(nodePath)

	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateFile(nodePath+item.GetName(), string(yaml))
	}
}

func SaveStorageClasses(clientset *kubernetes.Clientset, codec runtime.Codec, wg *sync.WaitGroup) {
	defer wg.Done()
	kubeObject, err := clientset.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate storageclasses using the current config.")
	utilities.RecreateDir(storageclassPath)

	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateFile(storageclassPath+item.GetName(), string(yaml))
	}
}

func SaveClusterRoleBindings(clientset *kubernetes.Clientset, codec runtime.Codec, wg *sync.WaitGroup) {
	defer wg.Done()
	kubeObject, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate clusterrolebindings using the current config.")
	utilities.RecreateDir(clusterrolebindingPath)

	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateFile(clusterrolebindingPath+item.GetName(), string(yaml))
	}
}

func SaveClusterRoles(clientset *kubernetes.Clientset, codec runtime.Codec, wg *sync.WaitGroup) {
	defer wg.Done()
	kubeObject, err := clientset.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate clusterroles using the current config.")
	utilities.RecreateDir(clusterrolePath)

	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateFile(clusterrolePath+item.GetName(), string(yaml))
	}
}

func SavePersistentVolumes(clientset *kubernetes.Clientset, codec runtime.Codec, wg *sync.WaitGroup) {
	defer wg.Done()
	kubeObject, err := clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate persistent volumes using the current config.")
	utilities.RecreateDir(persistentVolumePath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateFile(persistentVolumePath+item.GetName(), string(yaml))
	}
}

func SaveNamespaces(clientset *kubernetes.Clientset, codec runtime.Codec) []string {
	var namespaces []string
	kubeObject, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate namespaces using the current config.")
	utilities.RecreateDir(namspacePath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		namespaces = append(namespaces, item.Name)
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateFile(namspacePath+item.GetName(), string(yaml))
	}
	return namespaces
}

func SaveConfigMaps(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate config maps using the current config.")
	utilities.RecreateDir(configMapPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(configMapPath + namespace)
		utilities.CreateFile(configMapPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SavePersistentVolumeClaims(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate persistent volume claims using the current config.")
	utilities.RecreateDir(persistentVolumeClaimPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(persistentVolumeClaimPath + namespace)
		utilities.CreateFile(persistentVolumeClaimPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SavePods(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate pods using the current config.")
	utilities.RecreateDir(podPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(podPath + namespace)
		utilities.CreateFile(podPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveSecrets(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate secrets using the current config.")
	utilities.RecreateDir(secretPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(secretPath + namespace)
		utilities.CreateFile(secretPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveServiceAccounts(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate service accounts using the current config.")
	utilities.RecreateDir(serviceAccountPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(serviceAccountPath + namespace)
		utilities.CreateFile(serviceAccountPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveServices(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate services using the current config.")
	utilities.RecreateDir(servicePath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(servicePath + namespace)
		utilities.CreateFile(servicePath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveDaemonsets(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate daemonsets using the current config.")
	utilities.RecreateDir(daemonsetPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(daemonsetPath + namespace)
		utilities.CreateFile(daemonsetPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveDeployments(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate deployments using the current config.")
	utilities.RecreateDir(deploymentPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(deploymentPath + namespace)
		utilities.CreateFile(deploymentPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveReplicaSets(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate replicasets using the current config.")
	utilities.RecreateDir(replicasetPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(replicasetPath + namespace)
		utilities.CreateFile(replicasetPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveStatefulSets(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate statefulsets using the current config.")
	utilities.RecreateDir(statefulsetPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(statefulsetPath + namespace)
		utilities.CreateFile(statefulsetPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveCronJobs(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate cronjobs using the current config.")
	utilities.RecreateDir(cronjobPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(cronjobPath + namespace)
		utilities.CreateFile(cronjobPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveJobs(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate jobs using the current config.")
	utilities.RecreateDir(jobPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(jobPath + namespace)
		utilities.CreateFile(jobPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveIngresses(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate ingresses using the current config.")
	utilities.RecreateDir(ingressPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(ingressPath + namespace)
		utilities.CreateFile(ingressPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveRoleBindings(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate rolebindings using the current config.")
	utilities.RecreateDir(rolebindingPath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(rolebindingPath + namespace)
		utilities.CreateFile(rolebindingPath+namespace+"/"+item.GetName(), string(yaml))
	}
}

func SaveRoles(clientset *kubernetes.Clientset, codec runtime.Codec, namespace string) {
	kubeObject, err := clientset.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{})
	utilities.CheckIfError(err, "Unable to iterate roles using the current config.")
	utilities.RecreateDir(rolePath)

	clientset.RESTClient().Get()
	for index := range kubeObject.Items {
		item := kubeObject.Items[index]
		yaml, err := runtime.Encode(codec, &item)
		utilities.CheckIfError(err, "Unable to encode: "+item.Name+".")
		utilities.CreateDir(rolePath + namespace)
		utilities.CreateFile(rolePath+namespace+"/"+item.GetName(), string(yaml))
	}
}
