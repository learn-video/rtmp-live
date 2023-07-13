package discovery

import (
	"log"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

func Watch() {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Create, watcher.Write)
	r := regexp.MustCompile(".m3u8$")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Println(event)
			case err := <-w.Error:
				log.Println(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("/hls"); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 200); err != nil {
		log.Fatalln(err)
	}
}
