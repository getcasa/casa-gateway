package gateway

import (
	"net/http"

	"github.com/labstack/echo"
)

// GetDevices route return all devices
func GetDevices(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// AddDevice route return all devices
func AddDevice(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
