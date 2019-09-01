package collector

import (
    "fmt"
    "os"
    "bufio"
    "time"
    "strings"
    "encoding/json"
    "github.com/gocolly/colly"
)

const MAX_USER_HISTORY = 2000

type UserComment struct {
    MovieName string  `json:"movie_name"`
    MovieID string    `json:"movie_id"`
    Score string      `json:"score"`
    Date string       `json:"date"`
}

type UserCommentList struct {
    Uid string                 `json:"user_id"`
    Uname string               `json:"user_name"`
    UserHistory []*UserComment `json:"user_comments"`
}

type DoubanUserHistoryCollector struct {
    UserCommentList
    DoubanColly *colly.Collector
}

var DoubanUserHistoryHandler = &DoubanUserHistoryCollector{}

func (dh *DoubanUserHistoryCollector) FetchUserHistory() error {
    dh.UserHistory = make([]*UserComment, 0, MAX_USER_HISTORY)
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
    dh.DoubanColly.OnHTML("div[class=grid-view] div[class=info]>ul",func(e *colly.HTMLElement){
        user_comment := &UserComment{}
        movie_url := e.ChildAttr("li>a[href]", "href")
        if movie_url != "" {
            tmp := strings.Split(movie_url, "/")
            user_comment.MovieID = tmp[len(tmp) - 2]
        } else {
            user_comment.MovieID = ""
        }
        user_comment.MovieName = e.ChildText("a[href]>em")
        user_comment.Score = e.ChildAttr("li>span[class^=rating]", "class")
        user_comment.Date = e.ChildText("li>span[class=date]")

        //fmt.Println("Find text:", item_name)
        dh.UserHistory = append(dh.UserHistory, user_comment)
    })
    dh.DoubanColly.OnHTML("div[class=paginator] span[class=next] a[href]", func(e *colly.HTMLElement){
        next_page := e.Attr("href")
        url_next := "https://movie.douban.com" + next_page
        fmt.Println("fetching url_next: ", url_next)
        if e.Request.Depth < 1000 {
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
    fileName := dh.Uid + ".txt"
    f, err := os.OpenFile("./save_result/user_history/" + fileName, os.O_CREATE|os.O_RDWR, 0666)
    defer f.Close()
    if err != nil {
        fmt.Println(err.Error())
    } else {
        file_content, err := json.Marshal(dh.UserCommentList)
        if err != nil {
            fmt.Println(err.Error())
        } else {
            newWriter := bufio.NewWriterSize(f, 1024)
            if _, err = newWriter.Write(file_content); err != nil {
                fmt.Println(err)
            }
            if err = newWriter.Flush(); err != nil {
                fmt.Println(err)
            }
            //fmt.Println(string(file_content))
        }
    }
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
    fmt.Println("url is ", "https://movie.douban.com")
    dh.DoubanColly.Visit("https://movie.douban.com")
    dh.DoubanColly.Wait()
    return nil
}
