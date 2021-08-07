package k8client

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var invalidNamespaces = []string{
	"default",
	"kube-node-lease",
	"kube-public",
	"kube-system",
}

func isInvalid(name string) error {
	for _, n := range invalidNamespaces {
		if n == name {
			return fmt.Errorf("invalid namespace")
		}
	}

	return nil
}

func (c *client) GetNamespace(name string) (ret *corev1.Namespace, err error) {
	return c.clientset.CoreV1().Namespaces().Get(c.cfg.Ctx, name, metav1.GetOptions{})
}

func (c *client) GetNamespaces() (ret []*corev1.Namespace, err error) {
	list, err := c.clientset.CoreV1().Namespaces().List(c.cfg.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, l := range list.Items {
		if err = isInvalid(l.Name); err == nil {
			ret = append(ret, &l)
		}
	}

	return ret, nil
}

func (c *client) CreateNamespace(name string) (*corev1.Namespace, error) {
	c.log("Attempting to create namespace '%s'", name)
	if err := isInvalid(name); err != nil {
		return nil, err
	}

	nmsp, _ := c.GetNamespace(name)

	if nmsp != nil && nmsp.Status.Phase == corev1.NamespaceTerminating {
		c.log("Delaying 10 seconds while the existing namespace is delete... then will proceed with the creation")
		time.Sleep(time.Second * 10) // Wait 5 seconds for the namespace to be deleted
	}

	req := &corev1.Namespace{}
	req.Name = name

	nmsp, err := c.clientset.CoreV1().Namespaces().Create(c.cfg.Ctx, req, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	c.log("Successfully created namespace '%s'", name)
	return nmsp, nil
}

func (c *client) DeleteNamespace(name string) error {
	c.log("Attempting to delete namespace '%s'", name)
	if err := isInvalid(name); err != nil {
		return err
	}

	nmsp, err := c.clientset.CoreV1().Namespaces().Get(c.cfg.Ctx, name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Do not bother... the namespace is already terminating
	if nmsp.Status.Phase == corev1.NamespaceTerminating {
		c.log("Successfully deleted namespace '%s'", name)
		return nil
	}

	if err := c.clientset.CoreV1().Namespaces().Delete(c.cfg.Ctx, name, *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}

	c.log("Successfully deleted namespace '%s'", name)
	return nil
}
