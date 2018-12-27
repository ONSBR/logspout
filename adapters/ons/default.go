package ons

import (
	"fmt"
	"os"
	"strings"

	"github.com/gliderlabs/logspout/router"
)

func init() {
	router.AdapterFactories.Register(NewONSAdapter, "ons")
}

// NewONSAdapter returns a configured ons.Adapter
func NewONSAdapter(route *router.Route) (router.LogAdapter, error) {
	f, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return nil, err
	}
	defaultExcludeContainers := map[string]bool{"/logspout": true, "/registry": true, "/postgres": true, "/rabbitmq": true, "/proxy": true, "/adminer": true, "/influxdb": true, "/chronograf": true, "/portainer": true, "/mongo": true, "/git-server": true}
	excludeContainers := os.Getenv("EXCLUDE_CONTAINERS")
	if excludeContainers != "" {
		containersToExclude := strings.Split(excludeContainers, ",")
		defaultExcludeContainers = make(map[string]bool)
		for _, container := range containersToExclude {
			defaultExcludeContainers[container] = true
		}
	}
	return &Adapter{
		file:              f,
		route:             route,
		excludeContainers: defaultExcludeContainers,
	}, nil
}

// Adapter is a simple adapter that streams log output to a connection without any templating
type Adapter struct {
	file              *os.File
	route             *router.Route
	excludeContainers map[string]bool
}

// Stream sends log data to a connection
func (a *Adapter) Stream(logstream chan *router.Message) {
	defer a.file.Close()
	for message := range logstream {
		_, ok := a.excludeContainers[message.Container.Name]
		if ok {
			continue
		}
		_, err := a.file.WriteString(fmt.Sprintf("container(%s)\tsource(%s)\tlogspout_time(%s)\tlog_data(%s)\n", message.Container.Name, message.Source, message.Time, message.Data))
		if err != nil {
			fmt.Println("!!", err)
		}
	}

}
