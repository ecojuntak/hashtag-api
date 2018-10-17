## Prerequisite
- Install Golang
- Install RabbitMQ Server
- Install SQLite Database

## Download dependencies
go get

## Run the migration
```
go run main.go migrate
```

## Run the REST Server
```
go run main.go start-rest
```

## Run the AMQP Server
```
go run main.go start-amqp
```