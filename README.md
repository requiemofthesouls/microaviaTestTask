# Time Protocol client-server app (RFC-868)


## Параметры запуска приложения
|Параметр|Описание|Значение по умолчанию|
|-----|------|-------|
|-tcp|Использовать протокол TCP|true|
|-udp|Использовать протокол UDP|true|
|-workers|Количество воркеров (только для сервера)|2|
|-timeout|Таймаут на запись|10|
|-port|Порт|8080|
|-h|Адрес хоста назначения (только для клиента)|""|
|-с|Запуск клиента|false|

### Run TCP server locally
```shell
go mod vendor
go run main.go -tcp -port=8081
```

### Run UDP server locally
```shell
go mod vendor
go run main.go -udp -port=8081
```

### Run TCP client locally
```shell
go mod vendor
go run main.go -c -tcp -port=8081
```

### Run UDP client locally
```shell
go mod vendor
go run main.go -c -udp -port=8081
```

### Run TCP & UDP servers
```shell
docker-compose up tcp_server udp_server
```

### Run TCP & UDP clients
```shell
docker-compose up tcp_client udp_client
```