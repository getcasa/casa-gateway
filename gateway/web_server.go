package gateway

import (
	"net/http"

	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/getcasa/sdk"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// MessageResponse define json response for API
type MessageResponse struct {
	Message string `json:"message"`
}

// Version use SemVer
var Version = "0.0.1"

// StartWebServer start echo server
func StartWebServer(port string) {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	// V1
	v1 := e.Group("/v1")

	v1.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, MessageResponse{
			Message: "Welcome to Casa Gateway v" + Version,
		})
	})

	v1.GET("/discover/:plugin", func(c echo.Context) error {
		var discoveredDevices []sdk.DiscoveredDevice
		plugin := c.Param("plugin")

		if PluginFromName(plugin) != nil && PluginFromName(plugin).Discover != nil {
			gatewayID := utils.GetIDFile()
			result := PluginFromName(plugin).Discover()
			for _, res := range result {
				res.GatewayID = gatewayID
				discoveredDevices = append(discoveredDevices, res)
			}
		}

		return c.JSON(http.StatusOK, discoveredDevices)
	})

	v1.GET("/configs", func(c echo.Context) error {
		var configs []sdk.Configuration
		for _, localPlugin := range LocalPlugins {
			configs = append(configs, *localPlugin.Config)
		}

		return c.JSON(http.StatusOK, configs)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
