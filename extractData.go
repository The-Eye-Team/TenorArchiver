package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/labstack/gommon/color"
)

func downloadData(fileName, URL, dataType string, worker *sync.WaitGroup) error {
	defer worker.Done()

	// Create the file
	out, err := os.Create(arguments.Output + "/" + dataType + "/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		go downloadData(fileName, URL, dataType, worker)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func extractData(URL, id string, wg *sync.WaitGroup) {
	defer wg.Done()

	var worker sync.WaitGroup
	var fileName, mp4URL, webmURL, gifURL string

	c := colly.NewCollector()
	extensions.RandomUserAgent(c)

	c.OnHTML("head", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("meta", func(_ int, el *colly.HTMLElement) bool {
			if strings.Contains(el.Attr("content"), "/mp4") {
				mp4URL = el.Attr("content")
				return false
			}
			return true
		})
	})

	c.OnHTML("head", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("meta", func(_ int, el *colly.HTMLElement) bool {
			if strings.Contains(el.Attr("content"), "/webm") {
				webmURL = el.Attr("content")
				return false
			}
			return true
		})
	})

	c.OnHTML("head", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("meta", func(_ int, el *colly.HTMLElement) bool {
			if strings.Contains(el.Attr("content"), "tenor.gif") {
				gifURL = el.Attr("content")
				return false
			}
			return true
		})
	})

	c.OnResponse(func(r *colly.Response) {
		fileName = r.Request.URL.Path[6:]
	})

	c.Visit(URL)

	if len(fileName) > 0 {
		if len(mp4URL) > 0 {
			worker.Add(1)
			go downloadData(fileName+".mp4", mp4URL, "MP4", &worker)
		}
		if len(webmURL) > 0 {
			worker.Add(1)
			go downloadData(fileName+".webm", webmURL, "WebM", &worker)
		}
		if len(gifURL) > 0 {
			worker.Add(1)
			go downloadData(fileName+".gif", gifURL, "GIF", &worker)
		}
		fmt.Println(checkPre +
			color.Yellow(" [") +
			color.Green(id) +
			color.Yellow("] ") +
			color.Green("Downloaded!"))
	} else if arguments.Verbose {
		fmt.Println(crossPre +
			color.Yellow(" [") +
			color.Red(id) +
			color.Yellow("] ") +
			color.Red("No data found at this URL!"))
	}

	worker.Wait()
}
