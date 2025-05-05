#!/usr/bin/env sh

podman build -t gom-sender -f Dockerfile.sender &
podman build -t gom-receiver -f Dockerfile.receiver &
podman build -t gom-ng -f Dockerfile-nginx &
wait
podman kube play --configmap configs.yaml --replace deploy.yaml
