package eureka

import (
	"math/rand"
	"strings"
	"time"
)

const (
	pathApps = "apps"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (e *eureka) selectServiceURL() string {
	urls := e.config.Client.ServiceUrls
	if len(urls) == 0 {
		log.Fatal("There are no Eureka Service URLs found")
	}
	return strings.TrimSuffix(urls[rand.Int()%len(urls)], "/")
}

func (e *eureka) buildUrl(paths ...string) string {
	return strings.Join(append([]string{e.selectServiceURL()}, paths...), "/")
}
