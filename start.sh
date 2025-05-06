#!/usr/bin/env sh

podman build -t gom-sender -f Containerfile.sender &
podman build -t gom-receiver -f Containerfile.receiver &
podman build -t gom-ng -f Containerfile.nginx &
wait
podman kube play --configmap configs.yaml --replace deploy.yaml
