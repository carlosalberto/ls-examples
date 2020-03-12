//
// example code to test lightstep-tracer-go
//
// usage:
//   LIGHTSTEP_ACCESS_TOKEN=${SECRET_TOKEN} \
//   LIGHTSTEP_COMPONENT_NAME=demo-client-go \
//   LIGHTSTEP_SERVICE_VERSION=0.1.8 \
//   go run client.go

package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
)

var lsToken = os.Getenv("LIGHTSTEP_ACCESS_TOKEN")
var lsHost = os.Getenv("LIGHTSTEP_HOST")
var lsPort = os.Getenv("LIGHTSTEP_PORT")
var lsSecure = os.Getenv("LIGHTSTEP_SECURE")
var targetURL = os.Getenv("TARGET_URL")

func initLightstepTracer() {
	port, err := strconv.Atoi(lsPort)
	if err != nil {
		port = 8360
	}
	plaintext := false
	if lsSecure == "0" {
		plaintext = true
	}
	componentName := os.Getenv("LIGHTSTEP_COMPONENT_NAME")
	if len(componentName) == 0 {
		componentName = "test-go-server"
	}
	serviceVersion := os.Getenv("LIGHTSTEP_SERVICE_VERSION")
	if len(serviceVersion) == 0 {
		serviceVersion = "0.0.0"
	}
	endpoint := lightstep.Endpoint{Host: lsHost, Port: port, Plaintext: plaintext}
	opentracing.InitGlobalTracer(lightstep.NewTracer(lightstep.Options{
		AccessToken: lsToken,
		Collector:   endpoint,
		UseHttp:     true,
		Tags: opentracing.Tags{
			"lightstep.component_name": componentName,
			"service.version":          serviceVersion,
		},
		SystemMetrics: lightstep.SystemMetricsOptions{
			Endpoint: endpoint,
		},
	}))
}

func main() {
	initLightstepTracer()
	if len(targetURL) == 0 {
		targetURL = "http://localhost:8081"
	}
	for {
		contentLength := rand.Intn(2048)
		url := fmt.Sprintf("%s/content/%d", targetURL, contentLength)
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Request to %s, got %d bytes\n", url, res.ContentLength)
		}
		time.Sleep(1 * time.Second)
	}

}
