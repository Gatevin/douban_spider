package history_collector

import (
    "fmt"
    "time"
    "github.com/gocolly/colly"
)

type DoubanUserHistoryCollector struct {
    MainUrl string
    UserHistory []string
    Uid string
    Uname string
    DoubanColly *colly.Collector
}

var DoubanUserHistoryHandler = &DoubanUserHistoryCollector{}

func (dh *DoubanUserHistoryCollector) FetchHistoryWithUser() error {
    dh.UserHistory = make([]string, 0, 2000)
    dh.DoubanColly = colly.NewCollector(
    )
    dh.DoubanColly.Limit(&colly.LimitRule{
		DomainGlob: "*",
		Parallelism: 2,
		RandomDelay: 3*time.Second,
	})
    dh.DoubanColly.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
    })
    dh.DoubanColly.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited: ", r.Request.URL)
		//fmt.Println("Result is: ", string(r.Body))
    })
    dh.DoubanColly.OnHTML("div[class=item] a[href]>em",func(e *colly.HTMLElement){
        item_name := e.Text
        fmt.Println("Find text:", item_name)
        dh.UserHistory = append(dh.UserHistory, item_name)
    })
    dh.DoubanColly.OnHTML("div[class=paginator] span[class=next] a[href]", func(e *colly.HTMLElement){
        next_page := e.Attr("href")
        url_next := "https://movie.douban.com" + next_page
        fmt.Println("fetching url_next: ", url_next)
        if e.Request.Depth < 10 {
            e.Request.Visit(url_next)
        }
    })
    dh.DoubanColly.OnError(func(r *colly.Response, e error) {
	    fmt.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError", e)
	    fmt.Println("Retrying url: ", r.Request.URL)
	    r.Request.Retry()
    })
    url := "https://movie.douban.com/people/"+ dh.Uid + "/collect"
    fmt.Println("fetching url: ", url)
    dh.DoubanColly.Visit(url)
    dh.DoubanColly.Wait()
    fmt.Println(dh.UserHistory)
	return nil
}


func (dh *DoubanUserHistoryCollector) Test() error {
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
