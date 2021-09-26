# Casa-Gateway

Casa-gateway is here to connect devices together and get all datas to send it to the [Casa server](https://github.com/ItsJimi/casa). After all, server send actions to the gateway to control devices. Gateway use a [plugins](https://github.com/getcasa?q=plugin) system to control devices. Developed in Golang, it works on arm64 boards (e.g. raspberryPi, nas) and all amd64.

## Build

To build with plugins, you need to compile it in plugin mode and build casa-gateway with plugins compiled lib.

### arm64 (nas synology)

```
sudo env CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -o casa-gateway *.go
```

### amd64

```
go build -o casa-gateway *.go
```

## Launch

- You need a launched casa-server to connect our casa-gateway (https://github.com/getcasa/casa)
- Set env variable 'CASA_SERVER_PORT', check casa-server to set good value (default 4353)
- Init gateway

```
./casa-gateway init
```

- Start gateway

```
./casa-gateway start
```
