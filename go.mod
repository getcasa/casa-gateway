module github.com/ItsJimi/casa-gateway

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/getcasa/sdk v0.0.0-20190603120433-a6275a3eed49
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.9 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/oklog/ulid/v2 v2.0.2
	github.com/pborman/getopt v0.0.0-20190409184431-ee0cd42419d3 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20190621222207-cc06ce4a13d4
)

replace github.com/getcasa/sdk v0.0.0-20190603120433-a6275a3eed49 => ../casa-sdk
