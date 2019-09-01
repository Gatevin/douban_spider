package main


import (
	"fmt"
	"douban_spider/history_collector"
)

func main() {
	fmt.Println("Hello word！")
	history_collector.DoubanUserHistoryHandler.MainUrl="https://movie.douban.com/"
	history_collector.DoubanUserHistoryHandler.Uid = "summermonica"
	history_collector.DoubanUserHistoryHandler.Uname= "井颯"
	history_collector.DoubanUserHistoryHandler.FetchHistoryWithUser()
}
