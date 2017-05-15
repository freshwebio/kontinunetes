package k8s

import (
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client defines a subset of a Kubernetes API server
// client providing functionality we need.
type Client struct {
	Clientset *kubernetes.Clientset
}

// NewInClusterClient deals with creating a new
// instance of a Kubernetes client.
func NewInClusterClient() (*Client, error) {
	// Let's create an in cluster config.
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{Clientset: clientset}, nil
}

// NewClient deals with creating
// a new kubernetes client instance from provided configuration.
func NewClient(configFile string) (*Client, error) {
	// Create our configuration from the provided file.
	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{Clientset: clientset}, nil
}

// Deployments retrieves the deployments matching the provided label selector
// and filtered down for those with container specs with the provided image.
// TODO: Implement a way to distinguish between containers in different registries with the same name for completely
// different projects.
// TODO: Add caching layer for deployments.
func (c *Client) Deployments(
	namespace string,
	containerImage string,
	selector labels.Selector,
) (*v1beta1.DeploymentList, error) {
	opts := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	deployments, err := c.Clientset.ExtensionsV1beta1().Deployments(namespace).List(opts)
	if err != nil {
		return nil, err
	}
	c.filterDeploymentsByImage(containerImage, deployments)
	return deployments, nil
}

// Gets all deployments that have atleast a single container that ends with
// or is equal to the provided container image.
// The provided container image is expected to be of the form vendor/repo:tag.
func (c *Client) filterDeploymentsByImage(containerImage string, deployments *v1beta1.DeploymentList) {
	filtered := []v1beta1.Deployment{}
	for _, deployment := range deployments.Items {
		containers := deployment.Spec.Template.Spec.Containers
		if c.hasContainerWithImage(containerImage, containers) {
			filtered = append(filtered, deployment)
		}
	}
	deployments.Items = filtered
}

// Determines whether at least one of the provided set of containers have the specified
// image.
func (c *Client) hasContainerWithImage(image string, containers []v1.Container) bool {
	hasContainerImage := false
	i := 0
	for !hasContainerImage && i < len(containers) {
		if strings.HasSuffix(containers[i].Image, image) || containers[i].Image == image {
			hasContainerImage = true
		} else {
			i++
		}
	}
	return hasContainerImage
}

// ReplicationControllers retrieves the replication controllers matching the provided label selector
// and filters to rcs with containers which use the provided container image.
// TODO: Add caching layer for replication controllers.
func (c *Client) ReplicationControllers(
	namespace string,
	containerImage string,
	selector labels.Selector,
) (*v1.ReplicationControllerList, error) {
	opts := v1.ListOptions{
		LabelSelector: selector.String(),
	}
	replicationControllers, err := c.Clientset.CoreV1().ReplicationControllers(namespace).List(opts)
	if err != nil {
		return nil, err
	}
	c.filterRCsByImage(containerImage, replicationControllers)
	return replicationControllers, nil
}

// Deals with filtering a set of Replication Controllers down to those which use the provided image in
// at least one container.
func (c *Client) filterRCsByImage(image string, replicationControllers *v1.ReplicationControllerList) {
	filtered := []v1.ReplicationController{}
	for _, rc := range replicationControllers.Items {
		containers := rc.Spec.Template.Spec.Containers
		if c.hasContainerWithImage(image, containers) {
			filtered = append(filtered, rc)
		}
	}
	replicationControllers.Items = filtered
}

// RedeployDeployment deals with scaling down, deleting and then creating the provided deployment
// in the given namespace.
func (c *Client) RedeployDeployment(namespace string, deployment *v1beta1.Deployment) error {
	// First scale the deployment down to 0, this removes the pod replica instances.
	deplScale, err := c.Clientset.ExtensionsV1beta1().
		Scales(namespace).
		Get("Deployment", deployment.GetName())
	if err != nil {
		return err
	}
	_, err = c.Clientset.ExtensionsV1beta1().
		Scales(namespace).
		Update("Deployment", deplScale)
	if err != nil {
		return err
	}
	opts := &v1.DeleteOptions{}
	err = c.Clientset.ExtensionsV1beta1().
		Deployments(namespace).
		Delete(deployment.GetName(), opts)
	if err != nil {
		return err
	}
	// Make sure the resource version is empty so the system can set it.
	deployment.ResourceVersion = ""
	_, err = c.Clientset.ExtensionsV1beta1().
		Deployments(namespace).
		Create(deployment)
	if err != nil {
		return err
	}
	return nil
}

// RedeployRC deals with deleting and then creating the provided replication controller
// in the given namespace.
func (c *Client) RedeployRC(namespace string, rc *v1.ReplicationController) error {
	// Like with deployments scale the rc down to 0, this removes the pod replica instances.
	// This might not be needed but due to not enough clarity on whether ReplicationControlller dependants
	// get cleaned up immediately we'll do this for now.
	rcScale, err := c.Clientset.ExtensionsV1beta1().
		Scales(namespace).
		Get("ReplicationController", rc.GetName())
	if err != nil {
		return err
	}
	_, err = c.Clientset.ExtensionsV1beta1().
		Scales(namespace).
		Update("ReplicationController", rcScale)
	if err != nil {
		return err
	}
	opts := &v1.DeleteOptions{}
	err = c.Clientset.CoreV1().
		ReplicationControllers(namespace).
		Delete(rc.GetName(), opts)
	if err != nil {
		return err
	}
	// Make sure the resource version is empty so the system can set it.
	rc.ResourceVersion = ""
	_, err = c.Clientset.CoreV1().
		ReplicationControllers(namespace).
		Create(rc)
	if err != nil {
		return err
	}
	return nil
}
