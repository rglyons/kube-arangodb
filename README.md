# ArangoDB Kubernetes Operator

[![Docker Pulls](https://img.shields.io/docker/pulls/arangodb/kube-arangodb.svg)](https://hub.docker.com/r/arangodb/kube-arangodb/)

ArangoDB Kubernetes Operator helps do run ArangoDB deployments
on Kubernetes clusters.

To get started, follow the Installation instructions below and/or
read the [tutorial](./docs/Manual/Tutorials/Kubernetes/README.md).

## State

The ArangoDB Kubernetes Operator is still in **development**.

Running ArangoDB deployments (single, active-failover or cluster)
is reasonably stable, and we're in the process of validating
production readiness of various Kubernetes platforms.

The feature set of the ArangoDB Kubernetes Operator is close to what
it is intended to be.

[Documentation](./docs/README.md)

### Production readiness state

| Platform             | Kubernetes version | ArangoDB version | State | Production ready | Remarks |
|----------------------|--------------------|------------------|-------|------------------|-|
| Google GKE           | 1.10               | >= 3.3.13        | Runs  | Yes              | Don't use micro nodes |
| Amazon EKS           | 1.10               | >= 3.3.13        | Runs  | No               |
| Amazon & Kops        | 1.10               | >= 3.3.13        | Runs  | No               |
| Azure AKS            | 1.10               | >= 3.3.13        | ?     | No               |
| OpenShift            | 1.10               | >= 3.3.13        | Runs  | No               |
| Pivotal PKS          | 1.10               | >= 3.3.13        | ?     | No               |
| Scaleway Kubernetes  | 1.10               | >= 3.3.13        | ?     | No               |
| Bare metal (kubeadm) | 1.10               | >= 3.3.13        | Runs  | No               |
| Minikube             | 1.10               | >= 3.3.13        | Runs  | Not intended     |
| Docker for Mac Edge  | 1.10               | >= 3.3.13        | Runs  | Not intended     |

## Installation of latest release using Helm

```bash
# The following will install the operator for `ArangoDeployment` &
# `ArangoDeplomentReplication` resources.
helm install https://github.com/arangodb/kube-arangodb/releases/download/0.3.0/kube-arangodb.tgz
# To use `ArangoLocalStorage`, also run
helm install https://github.com/arangodb/kube-arangodb/releases/download/0.3.0/kube-arangodb-storage.tgz
```

## Installation of latest release using Kubectl

```bash
kubectl apply -f https://raw.githubusercontent.com/arangodb/kube-arangodb/0.3.0/manifests/crd.yaml
kubectl apply -f https://raw.githubusercontent.com/arangodb/kube-arangodb/0.3.0/manifests/arango-deployment.yaml
# To use `ArangoLocalStorage`, also run
kubectl apply -f https://raw.githubusercontent.com/arangodb/kube-arangodb/0.3.0/manifests/arango-storage.yaml
# To use `ArangoDeploymentReplication`, also run
kubectl apply -f https://raw.githubusercontent.com/arangodb/kube-arangodb/0.3.0/manifests/arango-deployment-replication.yaml
```

## Building

```bash
DOCKERNAMESPACE=<your dockerhub account> make
kubectl apply -f manifests/crd.yaml
kubectl apply -f manifests/arango-deployment-dev.yaml
# To use `ArangoLocalStorage`, also run
kubectl apply -f manifests/arango-storage-dev.yaml
# To use `ArangoDeploymentReplication`, also run
kubectl apply -f manifests/arango-deployment-replication-dev.yaml
```
