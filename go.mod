module github.com/ItsJimi/casa-gateway

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/getcasa/sdk v0.0.0-20190923145410-20bbee062dc8
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/pborman/getopt v0.0.0-20190409184431-ee0cd42419d3 // indirect
	github.com/pulento/go-ssdp v0.0.0-20180514024734-4a0ed625a78b // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/crypto v0.0.0-20190923035154-9ee001bba392
	golang.org/x/net v0.0.0-20190923162816-aa69164e4478 // indirect
	golang.org/x/sys v0.0.0-20190924092210-98129a5cf4a0 // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/getcasa/sdk v0.0.0-20190923145410-20bbee062dc8 => ../casa-sdk

go 1.13
