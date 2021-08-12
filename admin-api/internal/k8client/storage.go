package k8client

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *client) CreateNamespaceStorage(name, namespace string) (pvcUID string, err error) {
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
	// Delay for 5 seconds do that the storage will be available
	time.Sleep(time.Second * 5)

	return
}

func (c *client) DeleteNamespaceStorage(name, namespace string) error {
	return c.clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(c.cfg.Ctx, name, *metav1.NewDeleteOptions(0))
}

func (c *client) CreateNamespaceStorageDeployment(storageName, namespace string) (storageDeploymentPodID string, err error) {
	sdName := fmt.Sprintf("%s-deployment", storageName)
	createStDply := new(appsv1.Deployment)
	createStDply.ObjectMeta = metav1.ObjectMeta{
		Name:      sdName,
		Namespace: namespace,
		Labels: map[string]string{
			"component": fmt.Sprintf("%s-deployment", sdName),
			"type":      "storage-deployment",
		},
	}
	createStDply.Spec = appsv1.DeploymentSpec{
		Replicas: setReplicas(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: createStDply.ObjectMeta.Labels,
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: createStDply.ObjectMeta,
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{Name: storageName, VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: storageName,
						},
					}},
					{Name: "host", VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{Path: "/var/run"},
					}},
				},
				Containers: []v1.Container{
					{
						Name:    sdName,
						Image:   storageDeploymentImage.String(),
						Command: []string{"sleep"},
						Args:    []string{"infinity"},
						VolumeMounts: []v1.VolumeMount{
							{MountPath: "/host", Name: storageName},
							{MountPath: "/var/run", Name: "host"},
						},
					},
				},
			},
		},
	}

	if _, err = c.clientset.AppsV1().Deployments(namespace).Create(c.cfg.Ctx, createStDply, metav1.CreateOptions{}); err != nil {
		return
	}

	watch, err := c.clientset.CoreV1().Pods(namespace).Watch(c.cfg.Ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	for event := range watch.ResultChan() {
		p, ok := event.Object.(*v1.Pod)
		if !ok {
			return "", fmt.Errorf("unable to get pod data")
		}
		c.log("Waiting to pods for storage deployment (%s) to come online Status: '%s'\n", sdName, p.Status.Phase)
		if p.Status.Phase == "Running" {
			break
		}
		time.Sleep(time.Second * 5) // Delay for 5 seconds while waiting for pod to finish
	}

	pods, err := c.clientset.CoreV1().Pods(namespace).List(c.cfg.Ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	storageDeploymentPodID = pods.Items[0].Name

	return
}
