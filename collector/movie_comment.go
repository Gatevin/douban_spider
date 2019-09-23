package collector

import (
    "fmt"
    "os"
    "bufio"
    "time"
    "strings"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/proxy"
    . "douban_spider/utils"
)

const MAX_SIZE_OF_REVIEWS = 3000
const MAX_SIZE_OF_SHORT_COMMENTS = 20000

type MovieComment struct {
    UserName string  `json:"user_name"`
    UserID string    `json:"user_id"`
    Score string     `json:"score"`
    Date string      `json:"date"`
}

type MovieCommentList struct {
    MovieID string                     `json:"movie_id"`
    MovieName string                   `json:"movie_name"`
    MovieShortComments []*MovieComment `json:"movie_short_comments"`
    MovieReviews []*MovieComment       `json:"movie_reviews"`
}

type DoubanMovieCommentCollector struct {
    MovieCommentList
    DoubanColly *colly.Collector
    DoubanUserName string
    DoubanPassword string
    UseAccount bool

    AnonymousIp *IP
    Anonymous bool
}

var DoubanMovieCommentHandler = &DoubanMovieCommentCollector{}

func (dh *DoubanMovieCommentCollector) UseAnonymousIp(ip *IP) error {
    //ip := utils.Ipmgr.GetAnonymousIp()
    if ip != nil {
        dh.AnonymousIp = ip
        dh.Anonymous = true
    }
    return nil
}

func (dh *DoubanMovieCommentCollector) CancelUseAnonymousIp() error {
    dh.Anonymous = false
    return nil
}

func (dh *DoubanMovieCommentCollector) CancelUseAccount() error {
    dh.UseAccount = false
    return nil
}

func (dh *DoubanMovieCommentCollector) UseDoubanAccount(user_name string, password string) error {
    dh.DoubanUserName = user_name
    dh.DoubanPassword = password
    dh.UseAccount = true
    return nil
}


// 爬虫规则
func (dh *DoubanMovieCommentCollector) ConfigCollyRule() error {
    dh.DoubanColly = colly.NewCollector(
    )
    dh.DoubanColly.Limit(&colly.LimitRule{
        DomainGlob: "*",
        Parallelism: 2,
        RandomDelay: 5*time.Second,
    })
    if dh.Anonymous == true {
        proxyStr := "http://" + dh.AnonymousIp.Address + ":" + dh.AnonymousIp.Port
        fmt.Println(proxyStr)
        rp, err := proxy.RoundRobinProxySwitcher(proxyStr)
        if err != nil {
            fmt.Println("Colly Proxy set error")
            return nil
        }
        dh.DoubanColly.SetProxyFunc(rp)
    }
    if dh.UseAccount == true {
        err := dh.DoubanColly.Post("https://accounts.douban.com/j/mobile/login/basic",
                                    map[string]string{
                                            "name": dh.DoubanUserName,
                                            "password": dh.DoubanPassword,
                                    })
        if err != nil {
            fmt.Println(err.Error())
        }
    }
    dh.DoubanColly.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting: ", r.URL)
        r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
    })
    dh.DoubanColly.OnResponse(func(r *colly.Response) {
        fmt.Println("Visited: ", r.Request.URL)
        //fmt.Println("Result is: ", string(r.Body))
    })

    //这一部分先取影评
    dh.DoubanColly.OnHTML("header[class=main-hd]",func(e *colly.HTMLElement){
        movie_comment := &MovieComment{}
        user_url := e.ChildAttr("a[href][class=name]", "href")
        if user_url != "" {
            tmp := strings.Split(user_url, "/")
            movie_comment.UserID = tmp[len(tmp) - 2]
        } else {
            movie_comment.UserID = ""
        }
        movie_comment.UserName = e.ChildText("a[href][class=name]")
        movie_comment.Score = e.ChildAttr("span[class][title]", "title")
        movie_comment.Date = e.ChildText("span[content][class=main-meta]")

        dh.MovieReviews = append(dh.MovieReviews, movie_comment)
    })
    //影评的下一页
    dh.DoubanColly.OnHTML("div[class=paginator] span[class=next] a[href]", func(e *colly.HTMLElement){
        next_page := e.Attr("href")
        url_next := "https://movie.douban.com/subject/" + dh.MovieID + "/reviews" + next_page
        fmt.Println("fetching url_next: ", url_next)
        if e.Request.Depth < 1000 {
            e.Request.Visit(url_next)
        }
    })

    //这一部分取短评
    dh.DoubanColly.OnHTML("div[class=comment]", func(e *colly.HTMLElement){
        movie_comment := &MovieComment{}
        user_url := e.ChildAttr("span[class=comment-info]>a[href]", "href")
        if user_url != "" {
            tmp := strings.Split(user_url, "/")
            movie_comment.UserID = tmp[len(tmp) - 2]
        } else {
            movie_comment.UserID = ""
        }
        movie_comment.UserName = e.ChildText("span[class=comment-info]>a")
        movie_comment.Score = e.ChildAttr("span[class=comment-info]>span[class$=rating][title]", "title")
        movie_comment.Date = e.ChildAttr("span[class=comment-info]>span[class^=comment-time][title]", "title")

        dh.MovieShortComments = append(dh.MovieShortComments, movie_comment)
    })
    //短评的下一页
    dh.DoubanColly.OnHTML("div[id=paginator][class=center] a[href][class=next]", func(e *colly.HTMLElement){
        next_page := e.Attr("href")
        url_next := "https://movie.douban.com/subject/" + dh.MovieID + "/comments" + next_page
        fmt.Println("fetching url_next: ", url_next)
        if e.Request.Depth < 1000 {
            e.Request.Visit(url_next)
        }
    })

    dh.DoubanColly.OnError(func(r *colly.Response, e error) {
        fmt.Println("Request URL: ", r.Request.URL, " failed with error", e)
        //fmt.Println("Retrying url: ", r.Request.URL)
        //r.Request.Retry()
    })
    return nil
}

// 启动对一个作品id的爬虫
func (dh *DoubanMovieCommentCollector) FetchMovieComment() error {
    dh.MovieShortComments = make([]*MovieComment, 0, MAX_SIZE_OF_SHORT_COMMENTS)
    dh.MovieReviews = make([]*MovieComment, 0, MAX_SIZE_OF_REVIEWS)

    // sort=hotest 热度排序 sort=time 时间排序 rating=5 查看五星评分的
    review_url_star5 := "https://movie.douban.com/subject/"+ dh.MovieID + "/reviews" + "?rating=5"
    fmt.Println("fetching url: ", review_url_star5)
    dh.DoubanColly.Visit(review_url_star5)

    review_url_star4 := "https://movie.douban.com/subject/"+ dh.MovieID + "/reviews" + "?rating=4"
    fmt.Println("fetching url: ", review_url_star4)
    dh.DoubanColly.Visit(review_url_star4)


    // 热门 https://movie.douban.com/subject/26759819/comments?sort=new_score&status=P&percent_type=h
    // sort = time 最新 percent_type=h 表示好评(4，5星)
    short_comments_url_hot := "https://movie.douban.com/subject/" + dh.MovieID + "/comments" + "?sort=new_score&status=P&percent_type=h"
    fmt.Println("fetching url: ", short_comments_url_hot)
    dh.DoubanColly.Visit(short_comments_url_hot)

    short_comments_url_time := "https://movie.douban.com/subject/" + dh.MovieID + "/comments" + "?sort=time&status=P&percent_type=h"
    fmt.Println("fetching url: ", short_comments_url_time)
    dh.DoubanColly.Visit(short_comments_url_time)

    dh.DoubanColly.Wait()

    fileName := dh.MovieID + ".txt"
    f, err := os.OpenFile("./save_result/movie_comments/" + fileName, os.O_CREATE|os.O_RDWR, 0666)
    defer f.Close()
    if err != nil {
        fmt.Println(err.Error())
    } else {
        file_content, err := Json.Marshal(dh.MovieCommentList)
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

