# synthetic-service
![GitHub release (latest by date)](https://img.shields.io/github/v/release/checkelmann/synthetic-service)
[![Build Status](https://travis-ci.org/checkelmann/synthetic-service.svg?branch=master)](https://travis-ci.org/checkelmann/synthetic-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/checkelmann/synthetic-service)](https://goreportcard.com/report/github.com/checkelmann/synthetic-service)

This implements a synthetic-service for Keptn.

## Prerequirements
The keptn dynatrace-service must be installed in order to use the synthetic-service as it is using the Dynatrace credentials from it.

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

## Configuration

The Service will listen for the `sh.keptn.events.deployment-finished` event, and will create a Synthetic Monitor in Dynatrace with the `deploymentURIPubli` as check URL.
You can add a label `"SyntheticManuallyAssignedApp": "APPLICATION-XYZ` to your event to assign the monitor to an application.

## Development

Development can be conducted using any GoLang compatible IDE or Text-Editor (e.g., Jetbrains GoLand, VSCode with Go plugins).


## License

Please find more information in the [LICENSE](LICENSE) file.