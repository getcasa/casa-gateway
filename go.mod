module github.com/ItsJimi/casa-gateway

require (
	github.com/anvie/port-scanner v0.0.0-20180225151059-8159197d3770
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/getcasa/sdk v0.0.0-20191105095754-6df142bc28a9
	github.com/gorilla/websocket v1.4.1
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/valyala/fasttemplate v1.1.0 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20191105034135-c7e5f84aec59 // indirect
	golang.org/x/net v0.0.0-20191105084925-a882066a44e0 // indirect
	golang.org/x/sys v0.0.0-20191104094858-e8c54fb511f6 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20191104232314-dc038396d1f0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/getcasa/sdk v0.0.0-20190923145410-20bbee062dc8 => ../casa-sdk

go 1.13
