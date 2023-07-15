package discovery

import (
	"log"
	"net"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/radovskyb/watcher"
)

type Stream struct {
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
	Host     string `json:"host"`
}

func Watch(c Config) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create, watcher.Write)
	r := regexp.MustCompile(".m3u8$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				ReportStream(event.Path, c.IP)
			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(c.HLSPath); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 200); err != nil {
		log.Fatalln(err)
	}
}

func ReportStream(filename, ip string) {
	fullDir, file := filepath.Split(filename)
	dir := path.Base(fullDir)
	if ip != "" {
		ip = GetHostIP()
	}
	stream := Stream{Name: dir, Manifest: file, Host: ip}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(stream).
		Post("http://api:9090/streams")
	if err != nil {
		log.Print(err)
		return
	}
	log.Printf("Result from /streams API: %d", resp.StatusCode())
}

func GetHostIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
