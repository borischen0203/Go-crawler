package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)

type job struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Company string `json:"company"`
}

var (
	jobKey  map[string]string
	allJobs []job
)
var count = 0

// <a data-gnav-element-name="HiringLab" class="icl-GlobalFooter-link" href="https://www.hiringlab.org/uk/">Hiring Lab</a>
func main() {
	fmt.Println("[Crawler start]")
	// allJobs := make([]Job, 0)
	collector := colly.NewCollector(
		colly.AllowedDomains("uk.indeed.com"),
	)
	collector.Limit(&colly.LimitRule{
		// DomainRegexp: "",
		DomainGlob:  "uk.indeed.com/*",
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
		// Parallelism:  0,
	})

	jobKey = make(map[string]string)

	// collector.OnHTML(".job_seen_beacon .resultContent .jobTitle-color-purple > span", func(e *colly.HTMLElement) {
	// 	// collector.OnHTML(`.heading4 color-text-primary singleLineTitle tapItem-gutter`, func(e *colly.HTMLElement) {
	// 	// collector.OnHTML(".tapItem", func(e *colly.HTMLElement) {
	// 	// collector.OnHTML(".jobTitle-color-purple > span[title]", func(e *colly.HTMLElement) {
	// 	//
	// 	// goquerySelection := e.DOM
	// 	// fmt.Println(goquerySelection.Find("span").Children().Text())
	// 	// fmt.Println(e.Attr("title"))
	// 	// fmt.Println(e.Attr("title"))
	// 	// fmt.Println(e.Text)
	// 	fmt.Println(e.ChildText(".job_seen_beacon .resultContent .jobTitle-color-purple > span[title]"))
	// 	// count++
	// 	// fmt.Println(count)
	// })

	collector.OnHTML(".mosaic-provider-jobcards", func(e *colly.HTMLElement) {
		// if e.Attr("class") ==""
		// id := e.ChildText(".mosaic-provider-jobcards > a.tapItem > data-jk")
		// id := e.ChildText("a.data-jk")
		// id := e.ChildAttr("a", "data-jk")
		// id := e.ChildAttr(`a[id]`, "data-jk")
		// fmt.Println("eTEXT" + e.Text)
		// id := e.Attr("data-jk")
		// temp := job{}
		// // temp.ID = e.ChildAttr("a", "data-jk")
		// temp.Title = e.ChildAttr(".jobTitle-color-purple > span", "title")
		// count++

		e.ForEach(".slider_container", func(_ int, el *colly.HTMLElement) {
			temp := job{
				// ID:      el.ChildAttr(".tapItem", "data-jk"),
				Title:   el.ChildAttr(".jobTitle-color-purple > span", "title"),
				Company: el.ChildText("span.companyName"),
			}
			allJobs = append(allJobs, temp)
			fmt.Println(temp.ID + "," + temp.Title + "," + temp.Company)

		})
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for i := 1; i < 2; i++ {
		page := "&start=" + strconv.Itoa(i*10)
		link := fmt.Sprintf("https://uk.indeed.com/jobs?q=shop+assistant&l=London%s", page)
		collector.Visit(link)
	}

	writeJSON(allJobs)
}

func writeJSON(data []job) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("jobs.json", file, 0644)
}
