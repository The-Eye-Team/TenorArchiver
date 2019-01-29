package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/labstack/gommon/color"
)

var client = http.Client{}

var checkPre = color.Yellow("[") + color.Green("✓") + color.Yellow("]")
var tildPre = color.Yellow("[") + color.Green("~") + color.Yellow("]")
var crossPre = color.Yellow("[") + color.Red("✗") + color.Yellow("]")

func init() {
	// Disable HTTP/2: Empty TLSNextProto map
	client.Transport = http.DefaultTransport
	client.Transport.(*http.Transport).TLSNextProto =
		make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
}

func crawl() {
	var worker sync.WaitGroup
	var id string
	var count int

	// Loop through pages
	for index := arguments.StartID; index <= arguments.StopID; index++ {
		worker.Add(1)
		count++
		id = strconv.Itoa(index)
		go extractMP4("https://tenor.com/view/"+id, id, &worker)
		if count == arguments.Concurrency {
			worker.Wait()
			count = 0
		}
	}
}

func main() {
	// Parse arguments
	parseArgs(os.Args)

	// Create output directory
	os.MkdirAll(arguments.Output, os.ModePerm)

	crawl()
}
