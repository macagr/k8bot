package k8s

import (
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// NewClient establishes the connection to K3s (Local or In-Cluster)
func NewClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Not in cluster, falling back to local WSL kubeconfig")
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return kubernetes.NewForConfig(config)
}

func NewMetricsClient() (*metricsv.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return metricsv.NewForConfig(config)
}