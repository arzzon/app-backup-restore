package orchestratorClient

import (
	"github.com/arzzon/app-backup-restore/internal/constants"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

//var (
//	K8sClient *kubernetes.Clientset
//	once      sync.Once
//)

//func GetK8sClient() *kubernetes.Clientset {
//	// Create a Kubernetes client using the in-cluster configuration
//	once.Do(func() {
//		config, err := rest.InClusterConfig()
//		if err != nil {
//			panic(err.Error())
//		}
//		K8sClient, err = kubernetes.NewForConfig(config)
//		if err != nil {
//			panic(err.Error())
//		}
//	})
//	return K8sClient
//}

func GetClientFromKubeconfig(kubeconfigPath string) (*kubernetes.Clientset, error) {
	// Load kubeconfig file
	if kubeconfigPath == "" {
		homeDir := os.Getenv("HOME")
		kubeconfigPath = strings.Replace(constants.KUBECONFIG_PATH, "~", homeDir, 1)
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// Create Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
