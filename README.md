# README

This is a Keptn Service Template written in GoLang. 

Quick start:

1. Download this repo as a zip-file, extract it into a new folder named after the service you want to create (e.g., simple-service) 
1. Replace every occurrence of (docker) image names and tags from `checkelmann/synthetic-service` to your docker organization and image name (e.g., `yourorganization/simple-service`)
1. Replace every occurrence of `synthetic-service` with the name of your service (e.g., `simple-service`)
1. Push your code a Git repo (e.g., GitHub) and adapt all links that contain `github.com` (e.g., to `github.com/your-username/simple-service`)
1. Àdapt the [go.mod](go.mod) file and change `example.com/` to the actual package name (e.g., `github.com/your-username/simple-service`)
1. Add yourself to the [CODEOWNERS](CODEOWNERS) file
1. Initialize a git repository: 
  * `git init .`
  * `git add .`
  * `git commit -m "Initial Commit"`
1. Optional: Push to upstream git repo
1. Last but not least: Remove this intro within the README file 

# synthetic-service
![GitHub release (latest by date)](https://img.shields.io/github/v/release/checkelmann/synthetic-service)
[![Build Status](https://travis-ci.org/checkelmann/synthetic-service.svg?branch=master)](https://travis-ci.org/checkelmann/synthetic-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/checkelmann/synthetic-service)](https://goreportcard.com/report/github.com/checkelmann/synthetic-service)

This implements a synthetic-service for Keptn.

## Compatibility Matrix

| Keptn Version    | [synthetic-service Docker Image](https://hub.docker.com/r/checkelmann/synthetic-service/tags) |
|:----------------:|:----------------------------------------:|
|       0.6.1      | checkelmann/synthetic-service:0.1.0 |

## Installation

The *synthetic-service* can be installed as a part of [Keptn's uniform](https://keptn.sh).

### Deploy in your Kubernetes cluster

To deploy the current version of the *synthetic-service* in your Keptn Kubernetes cluster, apply the [`deploy/service.yaml`](deploy/service.yaml) file:

```console
kubectl apply -f deploy/service.yaml
```

This should install the `synthetic-service` together with a Keptn `distributor` into the `keptn` namespace, which you can verify using

```console
kubectl -n keptn get deployment synthetic-service -o wide
kubectl -n keptn get pods -l run=synthetic-service
```

### Up- or Downgrading

Adapt and use the following command in case you want to up- or downgrade your installed version (specified by the `$VERSION` placeholder):

```console
kubectl -n keptn set image deployment/synthetic-service synthetic-service=checkelmann/synthetic-service:$VERSION --record
```

### Uninstall

To delete a deployed *synthetic-service*, use the file `deploy/*.yaml` files from this repository and delete the Kubernetes resources:

```console
kubectl delete -f deploy/service.yaml
```

## Development

Development can be conducted using any GoLang compatible IDE or Text-Editor (e.g., Jetbrains GoLand, VSCode with Go plugins).

### Where to start

If you don't care about the details, your first entrypoint is [eventhandlers.go](eventhandlers.go). Within this file 
 you can add implementation for pre-defined Keptn Cloud events.
 
To better understand Keptn CloudEvents, please look at the [Keptn Spec](https://github.com/keptn/spec).
 
If you want to get more insights, please look into [main.go](main.go), [deploy/service.yaml](deploy/service.yaml),
 consult the [Keptn docs](https://keptn.sh/docs/) as well as existing [Keptn Core](https://github.com/keptn/keptn) and
 [Keptn Contrib](https://github.com/keptn-contrib/) services.

### Common tasks

* Build the binary: `go build -ldflags '-linkmode=external' -v -o synthetic-service`
* Run tests: `go test -race -v ./...`
* Build the docker image: `docker build . -t checkelmann/synthetic-service:dev` (Note: Replace `your-username` with your DockerHub account/organization)
* Push the docker image to DockerHub: `docker push checkelmann/synthetic-service:dev` (Note: Replace `your-username` with your DockerHub account/organization)
* Deploy the service using `kubectl`: `kubectl apply -f deploy/`
* Undeploy the service using `kubectl`: `kubectl deploy -f deploy/`
* Watch the deployment using `kubectl`: `kubectl -n keptn get deployment synthetic-service -o wide`
* Get logs using `kubectl`: `kubectl -n keptn logs deployment/synthetic-service -f`
* Watch the deployed pods using `kubectl`: `kubectl -n keptn get pods -l run=synthetic-service`
* Deploy the service using [Skaffold](https://skaffold.dev/): `skaffold run --tail` (Note: please adapt the image name in [skaffold.yaml](skaffold.yaml))

### Testing Cloud Events

We have dummy cloud-events in the form of PostMan Requests in the [test-events/](test-events/) directory.

## License

Please find more information in the [LICENSE](LICENSE) file.