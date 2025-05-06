package proto

//go:generate protoc account.proto --go_out=. --go-grpc_out=. --go-grpc_opt=paths=import --go_opt=paths=import
