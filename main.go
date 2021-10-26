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
	ID              string `json:"id"`
	Title           string `json:"title"`
	Company         string `json:"company"`
	CompanyLocation string `json:"companyLocation"`
	Salary          string `json:"Salary"`
}

var (
	jobKey  map[string]string
	allJobs []job
)
var count = 0

func main() {
	fmt.Println("[Crawler start]")

	//Set up collector
	collector := colly.NewCollector(
		colly.AllowedDomains("uk.indeed.com"),
		colly.Async(true),
	)

	// collector2 := collector.Clone()
	// collector3 := collector.Clone()

	collector.Limit(&colly.LimitRule{
		// DomainRegexp: "",
		DomainGlob:  "uk.indeed.com/*",
		Delay:       5 * time.Second,
		RandomDelay: 1 * time.Second,
		Parallelism: 5,
	})

	jobKey = make(map[string]string)

	collector.OnHTML(".mosaic-provider-jobcards", func(e *colly.HTMLElement) {

		e.ForEach(".slider_container", func(_ int, el *colly.HTMLElement) {
			temp := job{
				//TODO: optimize the goQuerySelector to get job id
				// ID:      el.ChildAttr(".tapItem", "data-jk"),
				Title:           el.ChildAttr(".jobTitle-color-purple > span", "title"),
				Company:         el.ChildText("span.companyName"),
				CompanyLocation: el.ChildText("div.companyLocation"),
				Salary:          el.ChildText(".salary-snippet-container .salary-snippet > span"),
			}

			//TODO: Use job id as key to check the jobs duplicate or not
			// if _, exist := jobKey[temp.ID]; !exist {
			// 	count++
			// 	jobKey[temp.ID] = ""
			// 	allJobs = append(allJobs, temp)
			// 	fmt.Println(temp.ID + "," + temp.Title + "," + temp.Company + "," + temp.CompanyLocation + "," + temp.Salary)
			// }

			allJobs = append(allJobs, temp)
			// fmt.Println(temp.ID + "," + temp.Title + "," + temp.Company + "," + temp.CompanyLocation + "," + temp.Salary)
			count++
		})
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
		fmt.Println("Visiting", r.URL.String())
	})

	for page := 0; page < 2; page++ {
		num := strconv.Itoa(page * 10)
		link := fmt.Sprintf("https://uk.indeed.com/jobs?q=shop+assistant&l=London&start=%s", num)
		collector.Visit(link)
	}
	// collector.Wait()
	// collector2.Wait()
	// collector3.Wait()

	writeJSON(allJobs)

	num := strconv.Itoa(count)
	fmt.Printf("Successfully crawled %s jobs", num)

}

func writeJSON(data []job) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}

	_ = ioutil.WriteFile("jobs.json", file, 0644)
}
