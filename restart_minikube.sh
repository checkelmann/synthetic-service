#!/bin/env bash
export CHANGE_MINIKUBE_NONE_USER=true
minikube stop
minikube delete
sudo -E minikube start --vm-driver=none --kubernetes-version v1.15.11 --cpus 2 --memory 2048
kubectl get nodes
echo y|keptn install --use-case=quality-gates --platform=kubernetes --gateway=NodePort --verbose
kubectl -n keptn set image deployment/bridge bridge=keptn/bridge2:20200326.0744 --record
kubectl apply -f dt_secret.yaml -n keptn
keptn create project sockshop --shipyard=./shipyard-quality-gates.yaml
keptn create service carts --project=sockshop

