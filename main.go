package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

type Job struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var count = 0

// <a data-gnav-element-name="HiringLab" class="icl-GlobalFooter-link" href="https://www.hiringlab.org/uk/">Hiring Lab</a>
func main() {
	fmt.Println("[Crawler start]")
	// allJobs := make([]Job, 0)
	collector := colly.NewCollector(
	// colly.AllowedDomains("uk.indeed.com"),
	)

	collector.OnHTML("span[title]", func(e *colly.HTMLElement) {
		// fmt.Println(e.Attr("title"))
		fmt.Println(e.Text)
		// count++
		// fmt.Println(count)
	})

	collector.Visit("https://uk.indeed.com/jobs?q=shop%20assistant&l=London")
}
