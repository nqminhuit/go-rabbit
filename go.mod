module server

go 1.24.2

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect

replace server/common => ./common

replace server/utils => ./utils
