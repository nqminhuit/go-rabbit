# go-rabbit
testing golang with rabbitmq

## Start RabbitMQ
```bash
podman run --replace -d --name rabbitmq -p 5672:5672 -p 15672:15672 docker.io/rabbitmq:4.0-management-alpine
```

- management dashboard is at: http://localhost:15672/
- credentials: `guest/guest`

## Start services
```bash
go run send.go
```


```bash
go run receiver/receive.go
```
