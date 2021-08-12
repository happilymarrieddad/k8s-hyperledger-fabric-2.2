package k8client

import (
	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"admin-api/types"
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

type Config struct {
	IsIncluster bool
	DebugLogs   bool
	Ctx         context.Context
}

type Client interface {
	// Certificate Authority
	CreateCertificateAuthority(
		namespace, persistentVolumeClaimName, username, password, scriptsMountPath,
		scriptsSubPath, pathToStoreCAFilesSubPath string, port int32, org *types.Organization,
	) (dplyName, dplyUID, svcName, svcUID, clientName, clientUID string, err error)
	DeleteCertificateAuthority(org *types.Organization, namespace string) error
	// Namespace
	GetNamespace(name string) (ret *corev1.Namespace, err error)
	GetNamespaces() (ret []*corev1.Namespace, err error)
	CreateNamespace(name string) (*corev1.Namespace, error)
	DeleteNamespace(name string) error
	// Storage
	CreateNamespaceStorage(name, namespace string) (pvcUID string, err error)
	DeleteNamespaceStorage(name, namespace string) error
	CreateNamespaceStorageDeployment(storageName, namespace string) (storageDeploymentPodID string, err error)
	// Pod Operations
	CopyFileToPod(srcPath string, destPath string, podName, containerName, namespace string) (err error)
}

func NewClient(cfg *Config) (Client, error) {
	newClient := new(client)
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}

	var err error
	if cfg.IsIncluster {
		// creates the in-cluster config
		newClient.restconfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		newClient.restconfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(newClient.restconfig)
	if err != nil {
		return nil, err
	}

	newClient.clientset = clientset
	newClient.cfg = cfg

	return newClient, nil
}

type client struct {
	clientset  *kubernetes.Clientset
	cfg        *Config
	restconfig *rest.Config
}

func (c *client) log(format string, a ...interface{}) {
	if c.cfg.DebugLogs {
		fmt.Printf(format+"\n", a...)
	}
}
