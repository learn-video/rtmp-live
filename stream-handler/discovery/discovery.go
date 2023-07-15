package discovery

import (
	"log"
	"path"
	"path/filepath"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

type Stream struct {
	Name     string `json:"name"`
	Manifest string `json:"manifest"`
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
	stream := Stream{Name: dir, Manifest: file}
	log.Println(stream)
}
