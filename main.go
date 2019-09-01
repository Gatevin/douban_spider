package main


import (
	"douban_spider/collector"
)

func main() {
	//collector.DoubanUserHistoryHandler.Uid = "summermonica"
	//collector.DoubanUserHistoryHandler.Uname= "井颯"
	//collector.DoubanUserHistoryHandler.FetchUserHistory()

    collector.DoubanMovieCommentHandler.MovieID = "26759819"
    collector.DoubanMovieCommentHandler.MovieName = "命运之夜——天之杯II ：迷失之蝶 劇場版Fate/stay night Heaven's Feel II.lost butterfly"
    collector.DoubanMovieCommentHandler.FetchMovieComment()
}
