module github.com/ItsJimi/casa-gateway

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fatih/color v1.7.0 // indirect
	github.com/getcasa/sdk v0.0.0-20190603120433-a6275a3eed49
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.9 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/crypto v0.0.0-20190618222545-ea8f1a30c443 // indirect
)

replace github.com/getcasa/sdk v0.0.0-20190603120433-a6275a3eed49 => ../casa-sdk
