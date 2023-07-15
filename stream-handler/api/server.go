package api

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

type Stream struct {
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
	Host     string `json:"host"`
}

type ReportHandler struct {
	rc *redis.Client
}

func NewReportHandler(rc *redis.Client) *ReportHandler {
	return &ReportHandler{rc: rc}
}

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

func (rh *ReportHandler) reportStream(c echo.Context) error {
	stream := new(Stream)
	if err := c.Bind(stream); err != nil {
		return err
	}
	ReportStream(stream, rh.rc)
	return nil
}

func RunServer() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	rc := NewRedis(cfg)
	rh := NewReportHandler(rc)
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/authorize", authorize)
	e.POST("/streams", rh.reportStream)
	log.Fatal(e.Start(":9090"))
}
