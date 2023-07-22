package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

type StreamData struct {
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
	if streamName == "golive" {
		log.Printf("Allowed ingest for: %s\n", streamName)
		locationURL := fmt.Sprintf("rtmp://127.0.0.1:1935/hls/%s", streamName)
		c.Response().Header().Set("Location", locationURL)
		return c.NoContent(302)
	}
	log.Printf("Ingest denied for: %s\n", streamName)
	return c.NoContent(500)
}

func (rh *ReportHandler) reportStream(c echo.Context) error {
	streamData := new(StreamData)
	if err := c.Bind(streamData); err != nil {
		return err
	}
	stream := &Stream{
		Name:     streamData.Name,
		Manifest: streamData.Manifest,
		Host:     streamData.Host,
	}
	if err := ReportStream(stream, rh.rc); err != nil {
		log.Print(err)
		c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func (rh *ReportHandler) fetchStream(c echo.Context) error {
	streamName := c.Param("stream")
	stream, err := FetchStream(streamName, rh.rc)
	if err == ErrStreamNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, stream)
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
	e.GET("/streams/:stream", rh.fetchStream)
	log.Fatal(e.Start(":9090"))
}
