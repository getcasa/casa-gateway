package gateway

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/oklog/ulid/v2"
)

// MessageResponse define json response for API
type MessageResponse struct {
	Message string `json:"message"`
}

// NewULID create an ulid
func NewULID() ulid.ULID {
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy)
}

// StartWebServer start echo server
func StartWebServer(port string) {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// V1
	v1 := e.Group("/v1")

	// Signup
	v1.POST("/signup", SignUp)

	// Devices
	v1.GET("/devices", GetDevices)
	v1.POST("/devices", AddDevice)

	e.Logger.Fatal(e.Start(":" + port))
}
