#!/usr/bin/env sh

podman build -t gom-sender -f Dockerfile.sender &
podman build -t gom-receiver -f Dockerfile.receiver &
wait
podman kube play --configmap configs.yaml --replace deploy.yaml
