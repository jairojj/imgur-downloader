package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("Missing url argument")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(os.Args))

	for _, link := range os.Args {
		go func(link string) {
			defer wg.Done()
			Download(link)
		}(link)
	}

	wg.Wait()

	fmt.Println("Finished")
}

//Download imgur img
func Download(url string) {
	c := colly.NewCollector()

	c.OnHTML("div.post-image img", func(e *colly.HTMLElement) {
		img := e.Attr("src")

		fmt.Println("Image found:", img)

		c.Visit(img)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			err := r.Save(r.FileName())
			if err != nil {
				panic(err)
			}
		}
	})

	c.Visit(url)
}
