<img src="https://raw.githubusercontent.com/scraly/gophers/main/men-in-black-v2.png" alt="men-in-black-v2" width=500>

<p align="Left">
  <p align="Left">
    <a href="https://github.com/borischen0203/go-crawler/actions/workflows/go.yml"><img alt="GitHub release" src="https://github.com/borischen0203/o-crawler/actions/workflows/go.yml/badge.svg?logo=github&style=flat-square"></a>
  </p>
</p>


# Go-crawler
This project mainly use for crawling job on Indeed.

# Features
- Scraping job on Indeed and output a json file

# How to use


## Run in Local:

Required
- Install go(version >= 1.6)

### Run process
Step1: Clone the repo
```bash
git clone https://github.com/borischen0203/go-crawpler.git
```
Step2: Run main file
```bash
go run main.go
```

### Custom
```bash
# Choose to scrap from first page to any page you want.
	for page := 0; page < 60; page++ {
		num := strconv.Itoa(page * 10)
		link := fmt.Sprintf("https://uk.indeed.com/jobs?q=shop+assistant&l=London&start=%s", num)
		collector.Visit(link)
	}
```

## Tech stack
- Golang
- Colly
- Html

## Todo:
- [ ] Should able to find the last page when crawling
- [ ] Should able to run for 30 second.
