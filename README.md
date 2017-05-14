# Kontinunetes

Kubernetes deployment webhook server with authentication.
This service can be run in or outside of a Kubernetes cluster.

## What's it for?
This application provides a HTTP server exposing endpoints as webhooks.
These webhooks are triggered by a container registry such as Dockerhub and deployments
and replication controllers tagged with the configured autoDeployLabel will be re-deployed.
The Kubernetes resources are identified by the name of the image in one of the spec containers in comparison to the image repo in the webhook payload.

For now this should only really be used for test environments and not in production
as it completely destroys and rebuilds a Deployment or ReplicationController stack on each push to a container registry.
This application expects the tag on the image push event and tag specified on one of the auto-deploy enabled resource containers to be the same.
Rolling updates are currently not supported and would recommended that rolling updates are handled as a manual deployment process for
production environments.

## Building application (Standalone)
To build this application standalone simply run `go build` from the root directory after ensuring
all the dependencies are current by ensuring you have the godep tool installed and running `godep restore`.

## Building application (Docker)
To build for docker you must firstly ensure all the dependencies are installed using `godep restore`.
Then run `CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .` to get a binary fully packaged with
all dependencies to run on the empty scratch base.
Now you can build the docker image and run it in docker as an out of cluster k8s client or in k8s
for in cluster usage.

## Configuration

The following configuration should be provided:

| Type   | Flag                          | Environment                    | File                          | Default value         |
| ------ | :---------------------------- |:------------------------------ |:----------------------------- | :-------------------- |
| string | -kubeconfig ./config          | KUBECONFIG="./config"          | kubeconfig ./config           | ""                    |
| string | -namespace myclstr            | NAMESPACE="myclstr"            | namespace myclstr             | "default"             |
| string | -apikey a653233fre343deweq1   | APIKEY="a653233fre343deweq1"   | apikey a653233fre343deweq1    | ""                    |
| string | -apikeyParamName apikey       | APIKEYPARAMNAME="apikey"       | apikeyParamName apikey        | "apikey"              |
| string | -autoDeployLabel auto-deploy  | AUTODEPLOYLABEL="autodeploy"   | autoDeployLabel autodeploy    | "autodeploy"          |

The apiKey parameter isn't required and if one isn't provided then no authentication occurs within this app.
Other routes would be to secure this webhook service with an API gateway. An example of this would be to use Kong.
You can then select from a wide variety of authentication methods and apply rate-limiting at the gateway level.

## Trigger support
Dockerhub is the only supported webhook trigger source so far.
