package gateway

import (
	"net/http"

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

	e.Logger.Fatal(e.Start(":" + port))
}
