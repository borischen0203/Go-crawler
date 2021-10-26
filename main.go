package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
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

func main() {
	fmt.Println("[Crawler start]")

	//Set up collector
	collector := colly.NewCollector(
		colly.AllowedDomains("uk.indeed.com"),
		colly.AllowURLRevisit(),
		colly.Async(true),
	)

	collector2 := collector.Clone()

	err := collector.Limit(&colly.LimitRule{
		// DomainRegexp: "",
		DomainGlob: `indeed\.com`,
		// Delay:       10 * time.Second,
		RandomDelay: 10 * time.Second,
		Parallelism: 2,
	})
	if err != nil {
		log.Fatal(err)
	}

	jobKey = make(map[string]string)

	collector.OnHTML("#mosaic-provider-jobcards", func(e *colly.HTMLElement) {

		e.ForEach(".tapItem", func(_ int, el *colly.HTMLElement) {
			crawlJob := job{
				ID:              el.Attr("data-jk"),
				Title:           el.ChildAttr(".jobTitle-color-purple > span", "title"),
				Company:         el.ChildText("span.companyName"),
				CompanyLocation: removeSuffix(el.ChildText("div.heading6.company_location.tapItem-gutter > pre > div")),
				Salary:          el.ChildText(".salary-snippet-container .salary-snippet > span"),
			}

			//Use job id as key to check the crawler job is duplicate or not
			if _, exist := jobKey[crawlJob.ID]; !exist { // if this job is not duplicate
				jobKey[crawlJob.ID] = ""
				allJobs = append(allJobs, crawlJob)
				fmt.Printf("ID: %s |Title:%s |Company: %s |Location: %s |Salary: %s \n", crawlJob.ID, crawlJob.Title, crawlJob.Company, crawlJob.CompanyLocation, crawlJob.Salary)
			} else {
				fmt.Println("-------------" + crawlJob.ID + " is duplicate -----------------")
			}
		})
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
		fmt.Println("Visiting", r.URL.String())
	})

	//TODO: allow to excute the program 10 second
	// p := context.TODO()
	// c, _ := context.WithTimeout(p, 5*time.Second)
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// start := time.Now()
	// go func(ctx context.Context) {
	// 	defer wg.Done()
	// 	for {
	// 		select {
	// 		case <-c.Done():
	// 			return
	// 		default:
	// 			for page := 0; page < 10; page++ {
	// 				num := strconv.Itoa(page * 10)
	// 				link := fmt.Sprintf("https://uk.indeed.com/jobs?q=shop+assistant&l=London&start=%s", num)
	// 				select {
	// 				case <-c.Done():
	// 					return
	// 				default:
	// 					collector.Visit(link)
	// 				}
	// 			}
	// 		}
	// 	}
	// }(c)
	// wg.Wait()
	// fmt.Println(time.Since(start))

	// TODO: Check is this last page
	for page := 0; page < 65; page++ {
		num := strconv.Itoa(page * 10)
		link := fmt.Sprintf("https://uk.indeed.com/jobs?q=shop+assistant&l=London&start=%s", num)
		collector.Visit(link)
	}
	collector.Wait()
	collector2.Wait()

	writeJSON(allJobs)
	fmt.Printf("Successfully crawled %d jobs", len(allJobs))
}

//Output all jobs to json file
func writeJSON(data []job) {
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("Unable to create json file")
		return
	}
	_ = ioutil.WriteFile("jobs.json", file, 0644)
}

//Remvoe Company location Suffix word
func removeSuffix(str string) string {
	if strings.ContainsAny(str, "+") {
		result := str[0:strings.IndexByte(str, '+')]
		return result
	} else {
		return str
	}
}
