package history_collector

import (
    "fmt"
    "time"
    "github.com/gocolly/colly"
)

type DoubanCollector struct {
    MainUrl string
    DoubanColly *colly.Collector
}

var DoubanHandler = &DoubanCollector{}

func (dh *DoubanCollector) Test() error {
    dh.DoubanColly = colly.NewCollector(
        //colly.AllowedDomains("https://movie.douban.com/"),
        //colly.Async(true),
    )
    dh.DoubanColly.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 2,
		RandomDelay: 5*time.Second,
	})
    dh.DoubanColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
    })
    dh.DoubanColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited: ", r.Request.URL)
		fmt.Println("Result is: ", string(r.Body))
    })
    dh.DoubanColly.OnError(func(r *colly.Response, e error) {
	    fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError", e)
	    fmt.Println("Retrying url: ", r.Request.URL)
	    r.Request.Retry()
    })
    fmt.Println("url is ", dh.MainUrl)
    dh.DoubanColly.Visit(dh.MainUrl)
    dh.DoubanColly.Wait()
    return nil
}
