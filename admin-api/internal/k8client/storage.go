package k8client

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *client) CreateStorage(name, namespace string) (pvcUID string, err error) {
	createPVC := &v1.PersistentVolumeClaim{}
	createPVC.ObjectMeta = metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
	qty, err := resource.ParseQuantity("2Gi")
	if err != nil {
		return
	}
	createPVC.Spec = v1.PersistentVolumeClaimSpec{
		AccessModes: []v1.PersistentVolumeAccessMode{v1.ReadWriteMany},
		Resources:   v1.ResourceRequirements{Requests: v1.ResourceList{"storage": qty}},
	}

	pvc, err := c.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(c.cfg.Ctx, createPVC, metav1.CreateOptions{})
	if err != nil {
		return
	}

	pvcUID = string(pvc.UID)

	return
}

func (c *client) DeleteStorge(name, namespace string) error {
	return c.clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(c.cfg.Ctx, name, *metav1.NewDeleteOptions(0))
}
