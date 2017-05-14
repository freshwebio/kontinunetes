package main

import (
	"github.com/freshwebio/kontinunetes/webhook"
	"github.com/julienschmidt/httprouter"
)

// Deals with registering our webhook routes.
// For now only dockerhub is supported.
// TODO: Support other docker image repository providers such as Quay.
func registerRoutes(router *httprouter.Router, ctrl *webhook.Controller) {
	router.POST("/auto-deploy/docker-hub", ctrl.AutoDeployDockerHub)
}
