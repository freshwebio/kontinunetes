package webhook

import (
	"encoding/json"
	"log"
	"net/http"

	"k8s.io/client-go/pkg/labels"
	"k8s.io/client-go/pkg/selection"

	"github.com/freshwebio/kontinunetes/k8s"
	"github.com/julienschmidt/httprouter"
)

// Controller provides us with the handler for webhook
// HTTP requests.
type Controller struct {
	k8sclient       *k8s.Client
	autoDeployLabel string
	namespace       string
}

// NewController instantiates a new controller instance
// for handling webhook requests.
func NewController(k8sclient *k8s.Client, autoDeployLabel, namespace string) *Controller {
	return &Controller{k8sclient: k8sclient, autoDeployLabel: autoDeployLabel, namespace: namespace}
}

// AutoDeployDockerHub handles dockerhub webhook requests.
func (c *Controller) AutoDeployDockerHub(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dhPayload := &DockerHubPayload{PushData: &DockerHubPushData{}, Repository: &DockerHubRepository{}}
	err := json.NewDecoder(r.Body).Decode(dhPayload)
	if err != nil {
		// Simply log error if we can't decode the provided data and go no futher.
		log.Println(err)
		return
	}
	image := dhPayload.Repository.RepoName + ":" + dhPayload.PushData.Tag
	selector := labels.NewSelector()
	req, err := labels.NewRequirement(c.autoDeployLabel, selection.Exists, []string{})
	if err != nil {
		log.Println(err)
		return
	}
	selector = selector.Add(*req)
	deployments, err := c.k8sclient.Deployments(c.namespace, image, selector)
	if err != nil {
		log.Println(err)
		return
	}
	replicationControllers, err := c.k8sclient.ReplicationControllers(c.namespace, image, selector)
	if err != nil {
		log.Println(err)
		return
	}
	// Now let's schedule a destroy and rebuild for each of our rcs and deployments.
	for _, deployment := range deployments.Items {
		err = c.k8sclient.RedeployDeployment(c.namespace, &deployment)
		if err != nil {
			// Log error and carry on.
			log.Println(err)
		}
	}
	for _, rc := range replicationControllers.Items {
		err = c.k8sclient.RedeployRC(c.namespace, &rc)
		if err != nil {
			log.Println(err)
		}
	}
	// Now simply let docker hub know we're ok.
	// This is most likely not needed but feels wrong even with a webhook to not provide a response
	// for a request over HTTP.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}
