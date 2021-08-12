package k8client

import (
	"admin-api/types"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateCertificateAuthority
// dplyName sets to "<org.NetworkName>-ca"
// scriptsMountPath default "/scripts"
// scriptsSubPath default "files/scripts"
// pathToStoreCAFilesSubPath default "state/<dplyName>"
// port default 7054
func (c *client) CreateCertificateAuthority(
	namespace, persistentVolumeClaimName, username, password, scriptsMountPath,
	scriptsSubPath, pathToStoreCAFilesSubPath string, port int32, org *types.Organization,
) (dplyName, dplyUID, svcName, svcUID, clientName, clientUID string, err error) {
	// TODO: need to ensure the files are available on the network

	dplyName = fmt.Sprintf("%s-ca", org.NetworkName)

	if len(scriptsMountPath) == 0 {
		scriptsMountPath = "/scripts"
	}
	if len(scriptsSubPath) == 0 {
		scriptsSubPath = "files/scripts"
	}
	if len(pathToStoreCAFilesSubPath) == 0 {
		pathToStoreCAFilesSubPath = fmt.Sprintf("state/%s", dplyName)
	}
	if port == 0 {
		port = 7054
	}

	// Create the Deployment
	caPathName := "/etc/hyperledger/fabric-ca-server"
	createDpl := new(appsv1.Deployment)
	createDpl.ObjectMeta = metav1.ObjectMeta{
		Name:      dplyName,
		Namespace: namespace,
		Labels: map[string]string{
			"organization": org.NetworkName,
			"component":    fmt.Sprintf("%s-deployment", org.NetworkName),
			"type":         "ca",
		},
	}
	createDpl.Spec = appsv1.DeploymentSpec{
		Replicas: setReplicas(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: createDpl.ObjectMeta.Labels,
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: createDpl.ObjectMeta,
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{Name: persistentVolumeClaimName, VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: persistentVolumeClaimName,
						},
					}},
				},
				Containers: []v1.Container{
					{
						Name:    dplyName,
						Image:   certificateAuthorityImage.String(),
						Command: []string{"sleep"},
						Args:    []string{"infinity"},
						Ports:   []v1.ContainerPort{{ContainerPort: port}},
						Env: []v1.EnvVar{
							{Name: "FABRIC_CA_HOME", Value: caPathName},
							{Name: "USERNAME", Value: username},
							{Name: "PASSWORD", Value: password},
							{Name: "CSR_HOSTS", Value: dplyName},
						},
						VolumeMounts: []v1.VolumeMount{
							{MountPath: scriptsMountPath, Name: persistentVolumeClaimName, SubPath: scriptsSubPath},
							{MountPath: caPathName, Name: persistentVolumeClaimName, SubPath: pathToStoreCAFilesSubPath},
						},
					},
				},
			},
		},
	}

	dpl, err := c.clientset.AppsV1().Deployments(namespace).Create(c.cfg.Ctx, createDpl, metav1.CreateOptions{})
	if err != nil {
		return
	}
	dplyUID = string(dpl.UID)

	// Create the Service
	svcName = fmt.Sprintf("%s-service", dplyName)
	createSvc := &v1.Service{}
	createSvc.ObjectMeta = metav1.ObjectMeta{
		Name:      svcName,
		Namespace: namespace,
		Labels: map[string]string{
			"organization": org.NetworkName,
			"component":    svcName,
			"type":         "ca",
		},
	}
	createSvc.Spec = v1.ServiceSpec{
		Type:     v1.ServiceTypeClusterIP,
		Selector: createSvc.ObjectMeta.Labels,
		Ports: []v1.ServicePort{
			{Port: port, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: port}},
		},
	}

	svc, err := c.clientset.CoreV1().Services(namespace).Create(c.cfg.Ctx, createSvc, metav1.CreateOptions{})
	if err != nil {
		return
	}
	svcUID = string(svc.UID)

	// Create the client
	clientName = fmt.Sprintf("%s-client", dplyName)
	createClient := new(appsv1.Deployment)
	createClient.ObjectMeta = metav1.ObjectMeta{
		Name:      clientName,
		Namespace: namespace,
		Labels: map[string]string{
			"organization": org.NetworkName,
			"component":    clientName,
			"type":         "ca",
		},
	}
	createClient.Spec = appsv1.DeploymentSpec{
		Replicas: setReplicas(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: createClient.ObjectMeta.Labels,
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: createClient.ObjectMeta,
			Spec: v1.PodSpec{
				Volumes: []v1.Volume{
					{Name: persistentVolumeClaimName, VolumeSource: v1.VolumeSource{
						PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
							ClaimName: persistentVolumeClaimName,
						},
					}},
				},
				Containers: []v1.Container{
					{
						Name:    clientName,
						Image:   certificateAuthorityImage.String(),
						Command: []string{"sleep", "infinity"},
						Env: []v1.EnvVar{
							{Name: "FABRIC_CA_SERVER_HOME", Value: "/etc/hyperledger/fabric-ca-client"},
							{Name: "ORG_NAME", Value: org.NetworkName},
							{Name: "CA_SCHEME", Value: "https"},
							{Name: "CA_URL", Value: "ibm-ca-service:7054"},
							{Name: "CA_USERNAME", Value: username},
							{Name: "CA_PASSWORD", Value: password},
							{Name: "CA_CERT_PATH", Value: "/etc/hyperledger/fabric-ca-server/tls-cert.pem"},
							{Name: "NUM_NODES", Value: "20"},
							{Name: "STARTING_INDEX", Value: "0"},
						},
						VolumeMounts: []v1.VolumeMount{
							{MountPath: scriptsMountPath, Name: persistentVolumeClaimName, SubPath: scriptsSubPath},
							{MountPath: "/state", Name: persistentVolumeClaimName, SubPath: "state"},
							{MountPath: "/files", Name: persistentVolumeClaimName, SubPath: "files"},
							{MountPath: caPathName, Name: persistentVolumeClaimName, SubPath: pathToStoreCAFilesSubPath},
							{MountPath: " /etc/hyperledger/fabric-ca-client", Name: persistentVolumeClaimName, SubPath: "state/ibm-ca-client"},
							{MountPath: "/etc/hyperledger/fabric-ca/crypto-config", Name: persistentVolumeClaimName, SubPath: "files/crypto-config"},
						},
					},
				},
			},
		},
	}

	client, err := c.clientset.AppsV1().Deployments(namespace).Create(c.cfg.Ctx, createClient, metav1.CreateOptions{})
	if err != nil {
		return
	}
	clientUID = string(client.UID)

	return
}

func (c *client) DeleteCertificateAuthority(org *types.Organization, namespace string) error {
	dplyName := fmt.Sprintf("%s-ca", org.NetworkName)
	if err := c.clientset.AppsV1beta1().Deployments(namespace).
		Delete(c.cfg.Ctx, dplyName, *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}

	if err := c.clientset.CoreV1().Services(namespace).
		Delete(c.cfg.Ctx, fmt.Sprintf("%s-service", dplyName), *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}

	if err := c.clientset.AppsV1beta1().Deployments(namespace).
		Delete(c.cfg.Ctx, fmt.Sprintf("%s-client", dplyName), *metav1.NewDeleteOptions(0)); err != nil {
		return err
	}

	return nil
}
