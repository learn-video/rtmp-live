package discovery

import (
	"log"
	"net"
	"path"
	"path/filepath"
	"regexp"
	"time"

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
				ReportStream(event.Path)
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

func ReportStream(filename string) {
	fullDir, file := filepath.Split(filename)
	dir := path.Base(fullDir)
	stream := Stream{Name: dir, Manifest: file, Host: GetHostIP()}
	log.Println(stream)
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
