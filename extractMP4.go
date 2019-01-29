package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/labstack/gommon/color"
)

func downloadMP4(fileName, mp4URL string) error {
	// Create the file
	out, err := os.Create(arguments.Output + "/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(mp4URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		downloadMP4(fileName, mp4URL)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func extractMP4(URL, id string, wg *sync.WaitGroup) {
	defer wg.Done()
	var fileName, mp4URL string
	c := colly.NewCollector()

	c.OnHTML("head", func(e *colly.HTMLElement) {
		e.ForEachWithBreak("meta", func(_ int, el *colly.HTMLElement) bool {
			if strings.Contains(el.Attr("content"), "/mp4") {
				mp4URL = el.Attr("content")
				return false
			}
			return true
		})
	})

	c.OnResponse(func(r *colly.Response) {
		fileName = r.Request.URL.Path[6:] + ".mp4"
	})

	c.Visit(URL)

	if len(fileName) > 0 && len(mp4URL) > 0 {
		downloadMP4(fileName, mp4URL)
		fmt.Println(checkPre +
			color.Yellow(" [") +
			color.Green(id) +
			color.Yellow("] ") +
			color.Green("Downloaded!"))
	} else {
		fmt.Println(crossPre +
			color.Yellow(" [") +
			color.Red(id) +
			color.Yellow("] ") +
			color.Red("No MP4 found at this URL!"))
	}
}
