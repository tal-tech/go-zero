package discovk8s

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"

	"time"
)

// K8sClient is the k8s regular client struct
type K8sClient struct {
	Clientset       *kubernetes.Clientset
	InformerFactory informers.SharedInformerFactory
	stop            chan struct{}
}

// NewK8sClient generate a kubernetes client
func NewK8sClient() (*K8sClient, error) {

	kubeconfig := getKubeconfig()

	// Get k8s config
	config, err := getK8sConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)

	stopCh := make(chan struct{})

	//factory.Start(stopCh)

	return &K8sClient{
		Clientset:       clientset,
		InformerFactory: factory,
		stop:            stopCh,
	}, nil
}

// getK8sConfig to get K8s config
func getK8sConfig(kubeconfig string) (*rest.Config, error) {
	var (
		config *rest.Config
		err    error
	)

	if kubeconfig != "" {
		// Use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		// Creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	config.Timeout = 10 * time.Second
	return config, nil
}

func getKubeconfig() string {

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	_, err := os.Lstat(kubeconfig)

	if err == nil || os.IsExist(err) {
		return kubeconfig
	}
	return ""
}
