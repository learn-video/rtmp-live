package auth

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func authorize(c echo.Context) error {
	streamName := c.QueryParam("name")
	log.Printf("Authorizing ingest for: %s\n", streamName)
	if streamName == "12345" {
		log.Printf("Allowed ingest for: %s\n", streamName)
		locationURL := fmt.Sprintf("rtmp://127.0.0.1:1935/hls/%s", streamName)
		c.Response().Header().Set("Location", locationURL)
		return c.NoContent(302)
	}
	log.Printf("Ingest denied for: %s\n", streamName)
	return c.NoContent(500)
}

func RunServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/authorize", authorize)
	log.Fatal(e.Start(":9090"))
}
