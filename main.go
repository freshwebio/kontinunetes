package main

import (
	"net/http"

	"github.com/freshwebio/entre"
	"github.com/freshwebio/kontinunetes/k8s"
	"github.com/freshwebio/kontinunetes/middleware"
	"github.com/freshwebio/kontinunetes/webhook"
	"github.com/julienschmidt/httprouter"
	"github.com/namsral/flag"
)

var (
	kubeconfig      = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	kubeNamespace   = flag.String("namespace", "default", "The namespace to use to update k8s containers in.")
	apiKey          = flag.String("apikey", "", "The API key to be used to authenticate requests to the server.")
	apiKeyParamName = flag.String("apikeyParamName", "apikey", "The name of the query parameter that is expected to hold the API key")
	autoDeployLabel = flag.String("autoDeployLabel", "kontinunetes.autodeploy", "The name the label to be used for selecting deployments and replication controllers to deploy")
)

func main() {
	flag.Parse()
	var err error
	var cli *k8s.Client
	if *kubeconfig == "" {
		// Let's create an in-cluster client.
		cli, err = k8s.NewInClusterClient()
		if err != nil {
			panic(err.Error())
		}
	} else {
		// If kube config flag is specified lets load our client
		// from config.
		cli, err = k8s.NewClient(*kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	// Get the webhook controller set up.
	ctrl := webhook.NewController(cli, *autoDeployLabel, *kubeNamespace)
	// Now setup our HTTP server.
	router := httprouter.New()
	registerRoutes(router, ctrl)
	e := entre.Basic()
	// Let's register our api key authentication middleware in the
	// case an API key is provided.
	if *apiKey != "" {
		e.Push(middleware.NewAuth(*apiKey, *apiKeyParamName))
	}
	e.PushHandler(router)
	http.ListenAndServe(":3229", e)
}
