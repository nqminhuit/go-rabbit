#!/usr/bin/env sh

podman build -t gom-sender -f Dockerfile.sender
podman build -t gom-receiver -f Dockerfile.receiver
podman kube play --replace deploy.yaml
