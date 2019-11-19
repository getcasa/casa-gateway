module github.com/ItsJimi/casa-gateway

require (
	github.com/anvie/port-scanner v0.0.0-20180225151059-8159197d3770
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/getcasa/sdk v0.0.0-20191119095609-3201367a4102
	github.com/gorilla/websocket v1.4.1
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/oklog/ulid/v2 v2.0.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/fasttemplate v1.1.0 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20191117063200-497ca9f6d64f // indirect
	golang.org/x/net v0.0.0-20191119073136-fc4aabc6c914 // indirect
	golang.org/x/sys v0.0.0-20191119060738-e882bf8e40c2 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20191119175705-11e13f1c3fd7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/getcasa/sdk v0.0.0-20191119095609-3201367a4102 => ../casa-sdk

go 1.13
