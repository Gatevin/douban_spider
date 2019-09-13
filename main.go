package main


import (
    "fmt"
    "douban_spider/config"
	"douban_spider/collector"
)

func main() {
    confVip := config.LoadConfig()
    useDoubanAccount := confVip.GetBool("doubanAccount.useAccount")
    if useDoubanAccount {
        userName := confVip.GetString("doubanAccount.userName")
        password := confVip.GetString("doubanAccount.password")
        fmt.Println("Douban accout using..\nAcount: ", userName, " Password: ", password)
        collector.DoubanMovieCommentHandler.UseDoubanAccount(userName, password)
    } else {
        fmt.Println("No douban account will be used")
    }

	//collector.DoubanUserHistoryHandler.Uid = "summermonica"
	//collector.DoubanUserHistoryHandler.Uname= "井颯"
	//collector.DoubanUserHistoryHandler.FetchUserHistory()
/*
    collector.DoubanMovieCommentHandler.MovieID = "26759819"
    collector.DoubanMovieCommentHandler.MovieName = "命运之夜——天之杯II ：迷失之蝶 劇場版Fate/stay night Heaven's Feel II.lost butterfly"
    collector.DoubanMovieCommentHandler.FetchMovieComment()
*/
}
