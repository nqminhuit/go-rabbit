# go-rabbit
testing golang with rabbitmq

## Generate code manually

pre-condition
```bash
sudo apt-get install -y protobuf-compiler
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

generate code
```bash
cd transport/grpc/proto/; protoc account.proto --go_out=. --go-grpc_out=. --go-grpc_opt=paths=import --go_opt=paths=import
```


## Start RabbitMQ
```bash
podman run --replace -d --name rabbitmq -p 5672:5672 -p 15672:15672 docker.io/rabbitmq:4.0-management-alpine
```

- management dashboard is at: http://localhost:15672/
- credentials: `guest/guest`

## Start individual service

sender:
```bash
SVPORT=9093 go run cmd/sender/main.go
```

receiver:
```bash
OPENSEARCH_INDEX_NAME_MDCORE=mdcorereports \
OPENSEARCH_USERNAME=admin \
OPENSEARCH_PASSWORD=5D27220@08e3 \
OPENSEARCH_ADDRESSES=https://localhost:9200 \
go run cmd/receiver/main.go
```

## Playground
The directory `_playground` is for POC only, to run code from that directory:
```bash
cd _playground
go run .
```

it will execute the main function in `main.go` file

## Full deployment
```bash
podman build -t gom-sender -f Dockerfile.sender
podman build -t gom-receiver -f Dockerfile.receiver
podman kube play --configmap configs.yaml --replace deploy.yaml
```

or:
```bash
bash start.sh
```
