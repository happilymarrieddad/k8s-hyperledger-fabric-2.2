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
	CreateStorage(name, namespace string) (pvcUID string, err error)
	DeleteStorge(name, namespace string) error
}

func NewClient(cfg *Config) (Client, error) {
	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}

	var (
		config *rest.Config
		err    error
	)

	if cfg.IsIncluster {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
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
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client{clientset, cfg}, nil
}

type client struct {
	clientset *kubernetes.Clientset
	cfg       *Config
}

func (c *client) log(format string, a ...interface{}) {
	if c.cfg.DebugLogs {
		fmt.Printf(format+"\n", a...)
	}
}
