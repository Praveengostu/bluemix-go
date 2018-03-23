package containerv1

import (
	"fmt"
	"github.com/IBM-Bluemix/bluemix-go/client"
)

type ClusterALB struct {
	ID                string      `json:"id"`
	Region            string      `json:"region"`
	DataCenter        string      `json:"dataCenter"`
	IsPaid            bool        `json:"isPaid"`
	IngressHostname   string      `json:"ingressHostname"`
	IngressSecretName string      `json:"ingressSecretName"`
	ALBs              []ALBConfig `json:"alb"`
}

// ALBConfig config for alb configuration
// swagger:model
type ALBConfig struct {
	AlbID          string `json:"albID" description:"The ALB id"`
	ClusterID      string `json:"clusterID"`
	Name           string `json:"name"`
	AlbType        string `json:"albType"`
	Enable         bool   `json:"enable" description:"Enable (true) or disable(false) ALB"`
	State          string `json:"state"`
	CreatedDate    string `json:"createdDate"`
	NumOfInstances string `json:"numOfInstances" description:"Desired number of ALB replicas"`
	Resize         bool   `json:"resize" description:"Indicate whether resizing should be done"`
	AlbIP          string `json:"albip" description:"BYOIP VIP to use for ALB. Currently supported only for private ALB"`
	Zone           string `json:"zone" description:"Zone to use for adding ALB. This is indicative of the AZ in which ALB will be deployed"`
}

// ALBTypes ...
type ALBTypes struct {
	Types []string `json:"types" description:"Types of ALB available in a region"`
}

// ClusterALBSecret albsecret related information for cluster
// swagger:model
type ClusterALBSecret struct {
	ID         string            `json:"id"`
	Region     string            `json:"region"`
	DataCenter string            `json:"dataCenter"`
	IsPaid     bool              `json:"isPaid"`
	ALBSecrets []ALBSecretConfig `json:"albSecrets" description:"All the ALB secrets created in this cluster"`
}

// ALBSecretConfig config for alb-secret configuration
// swagger:model
type ALBSecretConfig struct {
	SecretName          string `json:"secretName" description:"Name of the ALB secret"`
	ClusterID           string `json:"clusterID"`
	DomainName          string `json:"domainName" description:"Domain name of the certficate"`
	CloudCertInstanceId string `json:"cloudCertInstanceID" description:"Cloud Cert instance ID from which certficate is downloaded"`
	ClusterCrn          string `json:"clusterCrn"`
	CertCrn             string `json:"certCrn" description:"Unique CRN of the certficate which can be located in cloud cert instance"`
	IssuerName          string `json:"issuerName" description:"Issuer name of the certficate"`
	ExpiresOn           string `json:"expiresOn" description:"Expiry date of the certficate"`
	State               string `json:"state" description:"State of ALB secret"`
}

// ALBSecretsPerCRN ...
type ALBSecretsPerCRN struct {
	ALBSecrets []string `json:"albsecrets" description:"ALB secrets correponding to a CRN"`
}

//Clusters interface
type Albs interface {
	GetClusterALBs(clusterNameOrID string) ([]ALBConfig, error)
	GetClusterALB(albID string) (ALBConfig, error)
	ConfigureALB(albID string, config ALBConfig) error
	DeployALB(clusterNameOrID string, config ALBConfig) error
	DeployALBCert(config ALBSecretConfig) error
	UpdateALBCert(config ALBSecretConfig) error
	RemoveALBCert(clusterID string, secretName string, certCRN string) error
	GetALBCert(clusterID string, secretName string) (ALBSecretConfig, error)
	ListAllALBCerts(clusterID string) ([]ALBSecretConfig, error)
}

type alb struct {
	client *client.Client
}

func newAlbAPI(c *client.Client) Albs {
	return &alb{
		client: c,
	}
}

// GetClusterALBs returns the list of albs available for cluster
func (r *alb) GetClusterALBs(clusterNameOrID string) ([]ALBConfig, error) {
	var successV ClusterALB
	rawURL := fmt.Sprintf("/v1/alb/clusters/%s", clusterNameOrID)
	_, err := r.client.Get(rawURL, &successV)
	return successV.ALBs, err
}

// GetClusterALB returns details about particular alb for cluster
func (r *alb) GetClusterALB(albID string) (ALBConfig, error) {
	var successV ALBConfig
	_, err := r.client.Get(fmt.Sprintf("/v1/alb/albs/%s", albID), &successV)
	return successV, err
}

// ConfigureALB enables or disables alb for cluster
func (r *alb) ConfigureALB(albID string, config ALBConfig) error {
	var successV interface{}
	if config.Enable {
		_, err := r.client.Post("/v1/alb/albs", config, &successV)
		return err
	}
	_, err := r.client.Delete(fmt.Sprintf("/v1/alb/albs/%s", albID))
	return err
}

// ConfigureALB enables or disables alb for cluster
func (r *alb) DeployALB(clusterNameOrID string, config ALBConfig) error {
	var successV interface{}
	_, err := r.client.Put(fmt.Sprintf("/v1/alb/cluster/%s", clusterNameOrID), config, &successV)
	return err
}

// DeployALBCert deploys alb-cert
func (r *alb) DeployALBCert(config ALBSecretConfig) error {
	var successV interface{}
	_, err := r.client.Post("/v1/alb/albsecrets", config, &successV)
	return err
}

// UpdateALBCert updates alb-cert
func (r *alb) UpdateALBCert(config ALBSecretConfig) error {
	_, err := r.client.Put("/v1/alb/albsecrets", config, nil)
	return err
}

// RemoveALBCert removes the alb-cert
func (r *alb) RemoveALBCert(clusterID string, secretName string, certCRN string) error {
	var path string
	if secretName != "" {
		path = fmt.Sprintf("/v1/alb/clusters/%s/albsecrets?albSecretName=%s", clusterID, secretName)

	} else if certCRN != "" {
		path = fmt.Sprintf("/v1/alb/clusters/%s/albsecrets?certCrn=%s", clusterID, certCRN)
	}
	_, err := r.client.Delete(path)
	return err
}

// GetALBCert returns details about specified alb cert for given secretName
func (r *alb) GetALBCert(clusterID string, secretName string) (ALBSecretConfig, error) {
	var successV ALBSecretConfig
	_, err := r.client.Get(fmt.Sprintf("/v1/alb/clusters/%s/albsecrets?albSecretName=%s", clusterID, secretName), &successV)
	return successV, err
}

// ListAllALBCerts for cluster...
func (r *alb) ListAllALBCerts(clusterID string) ([]ALBSecretConfig, error) {
	var successV ClusterALBSecret
	_, err := r.client.Get(fmt.Sprintf("/v1/alb/clusters/%s/albsecrets", clusterID), &successV)
	return successV.ALBSecrets, err
}
