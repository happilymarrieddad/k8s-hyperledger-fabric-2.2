package k8client

type imageName string

func (i imageName) String() string {
	return string(i)
}

const (
	certificateAuthorityImage imageName = "hyperledger/fabric-ca:1.4.7"
	storageDeploymentImage    imageName = "nginx:stable"
)
